# Generated by Django 3.1.14 on 2022-03-25 01:54

from django.db import migrations


class Migration(migrations.Migration):
    dependencies = [
        ('api', '0010_auto_20220203_1055'),
    ]

    operations = [
        migrations.AlterModelTable(
            name='activeuser',
            table='active_user',
        ),
        migrations.AlterModelTable(
            name='division',
            table='division',
        ),
        migrations.AlterModelTable(
            name='floor',
            table='floor',
        ),
        migrations.AlterModelTable(
            name='hole',
            table='hole',
        ),
        migrations.AlterModelTable(
            name='message',
            table='message',
        ),
        migrations.AlterModelTable(
            name='pushtoken',
            table='push_token',
        ),
        migrations.AlterModelTable(
            name='report',
            table='report',
        ),
        migrations.AlterModelTable(
            name='tag',
            table='tag',
        ),
        migrations.AlterModelTable(
            name='user',
            table='user',
        ),
    ]
