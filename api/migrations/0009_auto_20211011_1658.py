# Generated by Django 3.1.13 on 2021-10-11 08:58

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('api', '0008_auto_20211011_1655'),
    ]

    operations = [
        migrations.RenameField(
            model_name='user',
            old_name='email',
            new_name='email_encrypted',
        ),
    ]
