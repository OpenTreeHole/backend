# Generated by Django 3.1.14 on 2022-01-19 07:13

from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ('api', '0002_floor_special_tag'),
    ]

    operations = [
        migrations.CreateModel(
            name='OldUserFavorites',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('uid', models.CharField(max_length=11)),
                ('favorites', models.JSONField()),
            ],
        ),
    ]
