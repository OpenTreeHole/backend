import random
from datetime import datetime, timezone

from django.conf import settings
from django.contrib.auth import get_user_model
from django.db.models import Case, When
from rest_framework import serializers

from api.models import Division, Tag, Hole, Floor, Report, Message
from utils.decorators import cache_function_call

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
    division_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = Division
        fields = ['division_id', 'name', 'description', 'pinned']

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
    tag_id = serializers.IntegerField(source='id', read_only=True, required=False)
    name = serializers.CharField(max_length=settings.MAX_TAG_LENGTH)
    temperature = serializers.IntegerField(required=False)

    class Meta:
        model = Tag
        fields = ['tag_id', 'name', 'temperature']


class FloorGetSerializer(serializers.Serializer):
    def create(self, validated_data):
        pass

    def update(self, instance, validated_data):
        pass

    hole_id = serializers.IntegerField(required=False, write_only=True, default=1)
    s = serializers.CharField(required=False, write_only=True)
    length = serializers.IntegerField(
        required=False, write_only=True,
        default=settings.PAGE_SIZE,
        max_value=settings.MAX_PAGE_SIZE,
        min_value=0
    )
    start_floor = serializers.IntegerField(
        required=False, write_only=True,
        default=0
    )


# 不序列化 mention 字段
class SimpleFloorSerializer(serializers.ModelSerializer):
    floor_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = Floor
        fields = ['floor_id', 'hole_id', 'content', 'anonyname', 'time_updated', 'time_created', 'deleted', 'fold', 'like']
        read_only_fields = ['floor_id', 'anonyname']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        return data

    @staticmethod
    def get_queryset(queryset):
        return queryset


class FloorSerializer(SimpleFloorSerializer):
    mention = SimpleFloorSerializer(many=True, read_only=True)

    class Meta:
        model = Floor
        fields = ['floor_id', 'hole_id', 'content', 'anonyname', 'mention', 'time_updated', 'time_created', 'deleted', 'fold', 'like']
        read_only_fields = ['floor_id', 'anonyname', 'mention']

    @staticmethod
    def get_queryset(queryset):
        return queryset.prefetch_related('mention')

    def to_representation(self, instance):
        # floor 使用缓存效果不好
        # @cache_function_call(f'floor#{instance.id}', settings.FLOOR_CACHE_SECONDS)
        def _inner_to_representation(self, instance):
            return super().to_representation(instance)

        data = _inner_to_representation(self, instance)
        user = self.context.get('user')
        if not user:
            print('[W] FloorSerializer 实例化时应提供参数 context={"user": request.user}')
        else:
            data['is_me'] = True if instance.user_id == user.id else False
            data['liked'] = True if user.id in instance.like_data else False
        return data


class MentionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Floor
        fields = ['mention']


class HoleSerializer(serializers.ModelSerializer):
    hole_id = serializers.IntegerField(source='id', read_only=True)
    division_id = serializers.IntegerField(default=1)
    tags = TagSerializer(many=True, required=False)
    length = serializers.IntegerField(
        required=False, write_only=True,
        default=settings.PAGE_SIZE,
        max_value=settings.MAX_PAGE_SIZE,
        min_value=1
    )
    prefetch_length = serializers.IntegerField(
        required=False, write_only=True,
        default=settings.FLOOR_PREFETCH_LENGTH,
        max_value=settings.MAX_PAGE_SIZE,
        min_value=1
    )
    start_time = serializers.DateTimeField(
        required=False, write_only=True,
        default=lambda: datetime.now(timezone.utc)  # 使用函数返回值，否则指向的是同一个对象
    )

    class Meta:
        model = Hole
        fields = ['hole_id', 'division_id', 'time_updated', 'time_created', 'tags', 'view', 'reply', 'length', 'prefetch_length', 'start_time']

    def to_representation(self, instance):
        """
        context 中传入 simple 字段，
            若为 True 则使用缓存并不返回所有与用户有关的数据
            若为 False 则不使用缓存，返回所有数据
        """

        def _inner_to_representation(self, instance):
            data = super().to_representation(instance)
            user = self.context.get('user')
            prefetch_length = self.context.get('prefetch_length', 1)
            if not user:
                print('[W] HoleSerializer 实例化时应提供参数 context={"user": request.user}')
            else:
                # serializer
                serializer = SimpleFloorSerializer if simple else FloorSerializer
                context = None if simple else {'user': user}

                # prefetch_data
                queryset = instance.floor_set.order_by('id')[:prefetch_length]
                queryset = serializer.get_queryset(queryset)
                prefetch_data = serializer(queryset, many=True, context=context).data

                # first_floor_data
                first_floor_data = prefetch_data[0] if len(prefetch_data) > 0 else None

                # last_floor_data
                queryset = instance.floor_set.order_by('-id')
                queryset = serializer.get_queryset(queryset)
                last_floor_data = serializer(queryset[0], context=context).data if len(queryset) > 0 else None

                data['floors'] = {
                    'first_floor': first_floor_data,
                    'last_floor': last_floor_data,
                    'prefetch': prefetch_data,
                }
            return data

        @cache_function_call(str(instance), settings.HOLE_CACHE_SECONDS)
        def _cached(self, instance):
            return _inner_to_representation(self, instance)

        simple = self.context.get('simple', False)
        if simple:
            return _cached(self, instance)
        else:
            return _inner_to_representation(self, instance)

    @staticmethod
    def get_queryset(queryset):
        return queryset.prefetch_related('tags')

    def validate_tags(self, tags):
        if len(tags) > settings.MAX_TAGS:
            raise serializers.ValidationError(f'标签不能多于{settings.MAX_TAGS}个', 400)
        return tags

    def validate_division_id(self, division_id):
        if not division_id:
            division, created = Division.objects.get_or_create(name='树洞')
            return division.pk
        elif not Division.objects.filter(pk=division_id).exists():
            raise serializers.ValidationError('分区不存在', 400)
        else:
            return division_id

    def create(self, validated_data):
        tags = validated_data.pop('tags')
        division_id = validated_data.pop('division_id')
        # 在添加 tag 前要保存 hole，否则没有id
        hole = Hole.objects.create(division_id=division_id)
        # 创建 tag 并添加至 hole
        for tag_name in tags:
            tag, created = Tag.objects.get_or_create(name=tag_name['name'])
            hole.tags.add(tag)
        return hole

    def update(self, instance, validated_data):
        tags = validated_data.get('tags')
        if tags:
            tag_list = []
            for tag_name in tags:
                tag, created = Tag.objects.get_or_create(name=tag_name['name'])
                tag_list.append(tag)
            instance.tags.set(tag_list)
        instance.view = validated_data.get('view', instance.view)
        instance.save()
        return instance


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


def add_a_floor(request, hole, category):
    """
    增加一条回复帖
    Args:
        request:
        hole:       hole对象
        category:   指定返回值为 floor 或 hole

    Returns:        floor or hole

    """
    # 校验 content
    serializer = FloorSerializer(data=request.data)
    serializer.is_valid(raise_exception=True)
    content = serializer.validated_data.get('content')
    # 校验 mention
    mention_serializer = MentionSerializer(data=request.data)
    mention_serializer.is_valid(raise_exception=True)
    mention = mention_serializer.validated_data.get('mention', [])

    # 获取匿名信息，如没有则随机选取一个，并判断有无重复
    anonyname = hole.mapping.get(str(request.user.pk))  # 存在数据库中的字典里的数据类型都是 string
    if not anonyname:
        while True:
            anonyname = random.choice(settings.NAME_LIST)
            if anonyname in hole.mapping.values():
                pass
            else:
                hole.mapping[request.user.pk] = anonyname
                break
    hole.save()

    # 创建 floor 并增加 hole 的楼层数
    floor = Floor.objects.create(hole=hole, content=content, anonyname=anonyname, user=request.user)
    floor.mention.set(mention)
    return hole if category == 'hole' else floor
