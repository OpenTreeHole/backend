# Generated by Django 3.2.7 on 2021-10-04 02:32

from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ('api', '0002_auto_20211004_0924'),
    ]

    operations = [
        migrations.AlterField(
            model_name='floor',
            name='mention',
            field=models.ManyToManyField(blank=True, related_name='mentioned_by', to='api.Floor'),
        ),
    ]