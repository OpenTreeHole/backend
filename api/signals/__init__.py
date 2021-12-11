from django.dispatch import Signal

modified_by_admin = Signal(providing_args=['instance'])
mention_to = Signal(providing_args=['instance', 'mentioned'])
new_penalty = Signal(providing_args=['instance', 'penalty'])
