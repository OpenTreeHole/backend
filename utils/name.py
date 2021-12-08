import json
import random

# filter names from txt list

# def select(filename, freq):
#     with open(filename, encoding='UTF-8') as f:
#         li = f.readlines()
#     a = filter(lambda i: int(i.split()[1]) >= freq, li)
#     b = map(lambda i: i.split()[0], a)
#     return list(b)
#
# li = select('food.txt', 1000) + select('animal.txt', 100) + select('reming.txt', 1000)

with open('utils/names.json', 'r', encoding='utf-8') as f:
    NAMES = json.load(f)

# with open('names.json', 'w', encoding='utf-8') as f:
#     json.dump(NAMES, f, ensure_ascii=False)

suffix = [
    '1', '2', '3', '4', '5'
]


def random_name(compare_set):
    cnt = 0
    while cnt < 100:
        name = random.choice(NAMES)
        if name not in compare_set:
            return name
        else:
            cnt += 1
    while True:
        name = random.choice(NAMES) + random.choice(suffix)
        if name not in compare_set:
            return name
        else:
            pass


if __name__ == '__main__':
    print(NAMES)
