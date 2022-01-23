# Generated by Django 3.1.14 on 2022-01-22 14:27

from django.db import migrations, models

import utils.default_values


class Migration(migrations.Migration):
    dependencies = [
        ('api', '0005_activeuser'),
    ]

    operations = [
        migrations.AlterField(
            model_name='activeuser',
            name='date',
            field=models.DateField(default=utils.default_values.default_active_user_date, unique=True),
        ),
    ]