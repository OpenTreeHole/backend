from rest_framework import serializers

from api.models import *


class UserSerializer(serializers.ModelSerializer):
    user_id = serializers.IntegerField(source='id')

    class Meta:
        model = User
        fields = ['user_id']


class ProfileSerializer(serializers.ModelSerializer):
    class Meta:
        model = Profile
        fields = ['user_id', 'favorites', 'permission']


class DivisionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Division
        fields = ['name', 'description']


class TagSerializer(serializers.ModelSerializer):
    class Meta:
        model = Tag
        fields = ['name', 'temperature']


class FloorSerializer(serializers.ModelSerializer):
    floor_id = serializers.IntegerField(source='id')

    class Meta:
        model = Floor
        fields = ['floor_id', 'content', 'anonyname', 'reply_to', 'time_updated', 'time_created', 'deleted', 'folded', 'like']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        user = self.context.get('user')
        if not user:
            print('[E] FloorSerializer 实例化时应提供参数 context={"user": request.user}')
        else:
            data['is_me'] = True if instance.user == user.pk else False
            data['liked'] = True if user.pk in instance.like_data else False
        return data


class HoleSerializer(serializers.ModelSerializer):
    hole_id = serializers.IntegerField(source='id')
    tags = TagSerializer(many=True)

    class Meta:
        model = Hole
        fields = ['hole_id', 'division_id', 'time_updated', 'time_created', 'tags', 'view', 'reply']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        user = self.context.get('user')
        if not user:
            print('[E] HoleSerializer 实例化时应提供参数 context={"user": request.user}')
        else:
            data['key_floor'] = {
                'first_floor': FloorSerializer(instance.floor_set.order_by('id')[0], context={'user': user}).data,
                'last_floor': FloorSerializer(instance.floor_set.order_by('-id')[0], context={'user': user}).data,
            }
        return data


class ReportSerializer(serializers.ModelSerializer):
    report_id = serializers.IntegerField(source='id')
    floor = FloorSerializer()

    class Meta:
        model = Report
        fields = ['report_id', 'hole_id', 'floor', 'reason', 'time_created', 'time_updated', 'dealed']

    def to_representation(self, instance):
        data = super().to_representation(instance)
        data['dealed_by'] = instance.dealed_by.profile.nickname if instance.dealed_by else None
        return data


class MessageSerializer(serializers.ModelSerializer):
    class Meta:
        model = Message
        fields = ['from_user', 'to_user', 'content', 'time_created']
