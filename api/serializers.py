from django.conf import settings
from django.contrib.auth import get_user_model
from django.db.models import Case, When
from rest_framework import serializers

from api.models import Division, Tag, Hole, Floor, Report, Message

User = get_user_model()


class UserSerializer(serializers.ModelSerializer):
    user_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = User
        fields = ['user_id', 'nickname', 'favorites', 'permission', 'config', 'joined_time', 'is_admin']

    def validate_permission(self, permission):
        for s in ['admin', 'silent']:
            if s not in permission:
                raise serializers.ValidationError(f'字段 {s} 不存在')
        return permission

    def validate_config(self, config):
        for s in ['show_folded', 'notify']:
            if s not in config:
                raise serializers.ValidationError(f'字段 {s} 不存在')
        return config


class DivisionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Division
        fields = ['name', 'description', 'pinned']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        order = Case(*[When(pk=pk, then=pos) for pos, pk in enumerate(instance.pinned)])  # Holes 按 pinned 的顺序排序
        holes_data = HoleSerializer(
            Hole.objects.filter(id__in=instance.pinned).order_by(order),
            many=True,
            context={'user': self.context.get('user')}
        ).data
        data['pinned'] = holes_data
        return data


class TagSerializer(serializers.ModelSerializer):
    tag_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = Tag
        fields = ['tag_id', 'name', 'temperature']


# noinspection DuplicatedCode
class FloorSerializer(serializers.ModelSerializer):
    class SubFloorSerializer(serializers.ModelSerializer):
        floor_id = serializers.IntegerField(source='id', read_only=True)

        class Meta:
            model = Floor
            fields = ['floor_id', 'hole_id', 'content', 'anonyname', 'mention', 'time_updated', 'time_created', 'deleted', 'fold', 'like']

        def to_representation(self, instance):
            data = super().to_representation(instance)
            user = self.context.get('user')
            if not user:
                print('[W] FloorSerializer 实例化时应提供参数 context={"user": request.user}')
            else:
                data['is_me'] = True if instance.user == user else False
                data['liked'] = True if user.pk in instance.like_data else False
            return data

    floor_id = serializers.IntegerField(source='id', read_only=True)
    mention = SubFloorSerializer(many=True, read_only=True)

    class Meta:
        model = Floor
        fields = ['floor_id', 'hole_id', 'content', 'anonyname', 'mention', 'time_updated', 'time_created', 'deleted', 'fold', 'like']
        read_only_fields = ['floor_id', 'anonyname', 'mention']

    def validate_content(self, content):
        content = content.strip()
        if not content:
            raise serializers.ValidationError('内容不能为空')
        return content

    def to_representation(self, instance):
        data = super().to_representation(instance)
        user = self.context.get('user')
        if not user:
            print('[W] FloorSerializer 实例化时应提供参数 context={"user": request.user}')
        else:
            data['is_me'] = True if instance.user == user else False
            data['liked'] = True if user.pk in instance.like_data else False
        return data


class HoleSerializer(serializers.ModelSerializer):
    hole_id = serializers.IntegerField(source='id', read_only=True)
    division_id = serializers.IntegerField(required=False)
    tags = TagSerializer(many=True, read_only=True)
    tag_names = serializers.ListField(required=False, write_only=True)

    class Meta:
        model = Hole
        fields = ['hole_id', 'division_id', 'time_updated', 'time_created', 'tags', 'tag_names', 'view', 'reply']

    def validate_tag_names(self, tag_names):
        if not tag_names:
            tag_names = ['默认']
        if len(tag_names) > settings.MAX_TAGS:
            raise serializers.ValidationError(f'标签不能多于{settings.MAX_TAGS}个', 400)
        for tag_name in tag_names:
            tag_name = tag_name.strip()
            if not tag_name:
                raise serializers.ValidationError('标签名不能为空', 400)
            if len(tag_name) > settings.MAX_TAG_LENGTH:
                raise serializers.ValidationError(f'标签名不能超过{settings.MAX_TAG_LENGTH}个字符', 400)
        return tag_names

    def validate_division_id(self, division_id):
        if not division_id:
            division, created = Division.objects.get_or_create(name='树洞')
            return division.pk
        elif not Division.objects.filter(pk=division_id).exists():
            raise serializers.ValidationError('分区不存在', 400)
        else:
            return division_id

    def to_representation(self, instance):
        data = super().to_representation(instance)
        user = self.context.get('user')
        if not user:
            print('[W] HoleSerializer 实例化时应提供参数 context={"user": request.user}')
        else:
            data['key_floor'] = {
                'first_floor': FloorSerializer(instance.floor_set.order_by('id')[0], context={'user': user}).data,
                'last_floor': FloorSerializer(instance.floor_set.order_by('-id')[0], context={'user': user}).data,
            }
        return data


class ReportSerializer(serializers.ModelSerializer):
    report_id = serializers.IntegerField(source='id', read_only=True)
    floor = FloorSerializer()

    class Meta:
        model = Report
        fields = ['report_id', 'hole_id', 'floor', 'reason', 'time_created', 'time_updated', 'dealed']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        data['dealed_by'] = instance.dealed_by.nickname if instance.dealed_by else None
        return data


class MessageSerializer(serializers.ModelSerializer):
    message_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = Message
        fields = ['message_id', 'message', 'code', 'data', 'has_read', 'time_created']
