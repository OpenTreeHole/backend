from datetime import datetime

from django.conf import settings
from django.contrib.auth import get_user_model
from django.contrib.auth.password_validation import validate_password
from django.core.cache import cache
from django.db.models import Case, When
from rest_framework import serializers

from api.models import Division, Tag, Hole, Floor, Report, Message, PushToken, ActiveUser
from api.signals import mention_to
from utils.apis import find_mentions
from utils.auth import many_hashes
from utils.decorators import cache_function_call
from utils.default_values import now
from utils.exception import BadRequest, Forbidden, ServerError
from utils.name import random_name

User = get_user_model()


class UserSerializer(serializers.ModelSerializer):
    user_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = User
        fields = ['user_id', 'nickname', 'favorites', 'permission', 'config',
                  'joined_time', 'is_admin']

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


class PushTokenSerializer(serializers.ModelSerializer):
    class Meta:
        model = PushToken
        fields = ['service', 'device_id', 'token']

    def validate_service(self, service):
        li = ['apns', 'mipush']
        if service not in li:
            raise serializers.ValidationError(f'字段需在 {li} 中')
        return service

    def create(self, validated_data):
        return PushToken.objects.create(**validated_data)


class BaseEmailSerializer(serializers.Serializer):
    email = serializers.EmailField()
    uuid = serializers.CharField(required=False)

    def create(self, validated_data):
        pass

    def update(self, instance, validated_data):
        pass


class EmailSerializer(BaseEmailSerializer):

    def validate_email(self, email):
        domain = email[email.find("@") + 1:]
        # 检查邮箱是否在白名单内
        if domain not in settings.EMAIL_WHITELIST:
            raise serializers.ValidationError('邮箱不在白名单内')
        return email


class RegisterSerializer(EmailSerializer):
    password = serializers.CharField()
    verification = serializers.CharField(max_length=6, min_length=6)

    def validate_password(self, password):
        validate_password(password)
        return password

    def validate(self, data):
        email = data['email']
        verification = data['verification']
        if not cache.get(email) or not cache.get(email) == verification:
            raise serializers.ValidationError('验证码错误')
        return data

    def create(self, validated_data):
        email = validated_data.get('email')
        password = validated_data.get('password')
        # 校验用户名是否已存在
        if User.objects.filter(identifier=many_hashes(email)).exists():
            raise BadRequest(detail='该用户已注册！如果忘记密码，请使用忘记密码功能找回')
        user = User.objects.create_user(email=email, password=password)
        cache.delete(email)  # 注册成功后验证码失效
        return user

    def update(self, instance, validated_data):
        instance.set_password(validated_data.get('password'))
        instance.save()
        return instance


class DivisionSerializer(serializers.ModelSerializer):
    division_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = Division
        fields = ['division_id', 'name', 'description', 'pinned']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        order = Case(
            *[When(pk=pk, then=pos) for pos, pk in enumerate(instance.pinned)]
        )  # Holes 按 pinned 的顺序排序
        holes_data = HoleSerializer(
            Hole.objects.filter(id__in=instance.pinned).order_by(order),
            many=True,
            context={'user': self.context.get('user')}
        ).data
        data['pinned'] = holes_data
        return data

    def update(self, instance, validated_data):
        instance.name = validated_data.get('name', instance.name)
        instance.description = validated_data.get('description', instance.description)
        instance.pinned = validated_data.get('pinned', instance.pinned)
        instance.save()
        return instance


class TagSerializer(serializers.ModelSerializer):
    tag_id = serializers.IntegerField(source='id', read_only=True, required=False)
    name = serializers.CharField(max_length=settings.MAX_TAG_LENGTH)
    temperature = serializers.IntegerField(required=False)

    class Meta:
        model = Tag
        fields = ['tag_id', 'name', 'temperature']

    def create(self, validated_data):
        tag, created = Tag.objects.get_or_create(name=validated_data.get('name'))
        if not created:
            raise BadRequest('tag 已存在')
        return tag

    def update(self, instance, validated_data):
        instance.name = validated_data.get('name', instance.name)
        instance.temperature = validated_data.get('temperature', instance.temperature)
        instance.save()
        return instance


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
    reverse = serializers.BooleanField(default=False)


# 不序列化 mention 字段
class SimpleFloorSerializer(serializers.ModelSerializer):
    floor_id = serializers.IntegerField(source='id', read_only=True)

    class Meta:
        model = Floor
        fields = ['floor_id', 'hole_id', 'content', 'anonyname', 'time_updated',
                  'time_created', 'deleted', 'fold', 'like', 'special_tag', 'storey']
        read_only_fields = ['floor_id', 'anonyname', 'storey']

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
        fields = ['floor_id', 'hole_id', 'content', 'history', 'anonyname', 'mention',
                  'time_updated', 'time_created', 'deleted', 'fold', 'like',
                  'special_tag', 'storey']
        read_only_fields = ['floor_id', 'history', 'anonyname', 'storey']

    @staticmethod
    def get_queryset(queryset):
        return queryset.prefetch_related('mention')

    def get_user(self):
        user = self.context.get('user')
        if not isinstance(user, get_user_model()):
            raise ServerError('FloorSerializer 实例化时应提供参数 context={"user": request.user}')
        return user

    def validate_special_tag(self, special_tag):
        user = self.get_user()
        if not user or not user.is_admin:
            raise Forbidden()
        return special_tag

    def validate_anonyname(self, anonyname):
        user = self.get_user()
        if not user.is_admin:
            raise Forbidden()
        return anonyname

    def create(self, validated_data):
        content = validated_data.get('content', '')
        special_tag = validated_data.get('special_tag', '')
        mentions = find_mentions(content)
        user = self.context.get('user')
        hole = self.context.get('hole')
        if not user or not hole:
            raise BadRequest(detail='创建floor需要在context中提供user和hole')
        # 获取匿名信息，如没有则随机选取一个，并判断有无重复
        anonyname = hole.mapping.get(str(user.pk))  # 存在数据库中的字典里的数据类型都是 string
        if not anonyname:
            anonyname = random_name(hole.mapping.values())
            hole.mapping[user.pk] = anonyname
        hole.reply += 1
        hole.save()
        floor = Floor.objects.create(hole=hole, content=content, anonyname=anonyname,
                                     user=user, special_tag=special_tag, storey=hole.reply)
        floor.mention.set(mentions)
        mention_to.send(sender=Floor, instance=floor, mentioned=mentions)
        return floor

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


class FloorUpdateSerializer(FloorSerializer):
    like = serializers.CharField(required=False)

    class Meta:
        model = Floor
        fields = ['content', 'anonyname', 'fold', 'like', 'special_tag']
        extra_kwargs = {
            'content': {'required': False},
            'anonyname': {'required': False}
        }


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
        default=now  # 使用函数返回值，否则指向的是同一个对象
    )

    class Meta:
        model = Hole
        fields = ['hole_id', 'division_id', 'time_updated', 'time_created', 'tags',
                  'view', 'reply', 'length', 'prefetch_length', 'start_time', 'hidden']

    def to_representation(self, instance):
        """
        context 中传入 simple 字段，
            若为 True 则使用缓存并不返回所有与用户有关的数据
            若为 False 则不使用缓存，返回所有数据
        """

        def _inner_to_representation(self, instance):
            data = super().to_representation(instance)
            user = self.context.get('user')
            prefetch_length = self.context.get('prefetch_length',
                                               settings.FLOOR_PREFETCH_LENGTH)
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
                queryset = serializer.get_queryset(instance.floor_set)
                last_floor_data = serializer(queryset.last(), context=context).data

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
        if len(tags) == 0:
            raise serializers.ValidationError('tags 不能为空', 400)
        if len(tags) > settings.MAX_TAGS:
            raise serializers.ValidationError(f'标签不能多于{settings.MAX_TAGS}个', 400)
        return tags

    def validate_division_id(self, division_id):
        @cache_function_call(division_id, 86400)
        def division_exists(division_id):
            return Division.objects.filter(pk=division_id).exists()

        if not division_exists(division_id):
            raise serializers.ValidationError('分区不存在', 400)
        else:
            return division_id

    def update(self, instance, validated_data):
        tags = validated_data.get('tags')
        if tags:
            tag_list = []
            for tag_name in tags:
                tag, created = Tag.objects.get_or_create(name=tag_name['name'])
                tag_list.append(tag)
            instance.tags.set(tag_list)
        instance.view = validated_data.get('view', instance.view)
        instance.division_id = validated_data.get('division_id', instance.division_id)
        instance.save()
        return instance


class HoleCreateSerializer(HoleSerializer):
    tags = TagSerializer(many=True)

    def create(self, validated_data):
        # 在添加外键前要保存 hole，否则没有id
        hole = Hole.objects.create(division_id=validated_data.get('division_id'))
        self.context.update({'hole': hole})
        floor_serializer = FloorSerializer(
            data=self.context.get('request_data'),
            context=self.context
        )
        floor_serializer.is_valid(raise_exception=True)
        floor_serializer.save()

        # 创建 tag 并添加至 hole
        for tag_name in validated_data.get('tags'):
            tag, created = Tag.objects.get_or_create(name=tag_name['name'])
            hole.tags.add(tag)
        self.context.get('user').favorites.add(hole)  # 自动收藏自己发的树洞
        return hole


class ReportSerializer(serializers.ModelSerializer):
    report_id = serializers.IntegerField(source='id', read_only=True)
    floor = FloorSerializer()

    class Meta:
        model = Report
        fields = ['report_id', 'hole_id', 'floor', 'reason', 'time_created',
                  'time_updated', 'dealed']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        data['dealed_by'] = instance.dealed_by.nickname if instance.dealed_by else None
        return data


class MessageSerializer(serializers.ModelSerializer):
    message_id = serializers.IntegerField(source='id', read_only=True)
    clear_all = serializers.BooleanField(default=False, write_only=True)
    not_read = serializers.BooleanField(default=True, write_only=True)
    start_time = serializers.DateTimeField(default=now, write_only=True)

    class Meta:
        model = Message
        fields = ['message_id', 'message', 'code', 'data', 'has_read', 'time_created',
                  'clear_all', 'not_read', 'start_time']

    def update(self, instance, validated_data):
        instance.message = validated_data.get('message', instance.message)
        instance.has_read = validated_data.get('has_read', instance.has_read)
        instance.code = validated_data.get('code', instance.code)
        instance.data = validated_data.get('data', instance.data)
        instance.save()
        return instance


class ActiveUserSerializer(serializers.ModelSerializer):
    start_date = serializers.DateField(default=datetime.now(settings.TIMEZONE).date(),
                                       write_only=True)
    end_date = serializers.DateField(default='1970-01-01', write_only=True)

    class Meta:
        model = ActiveUser
        fields = ['date', 'dau', 'mau', 'start_date', 'end_date']
