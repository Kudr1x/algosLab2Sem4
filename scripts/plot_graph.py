import matplotlib.pyplot as plt
import json
import os

# Путь к файлу с данными
json_path = '/home/kudrix/GolandProjects/AlgosSem4Lab2Neo/output/temp/compression_sizes.json'

try:
    with open(json_path, 'r') as f:
        file_sizes = json.load(f)
except FileNotFoundError:
    print(f"Ошибка: Файл {json_path} не найден!")
    exit(1)

# Определяем уровни качества из данных
expected_quality_levels = [1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100]

# Создаем директорию для графиков, если её нет
output_dir = 'graphs'
os.makedirs(output_dir, exist_ok=True)

# Создаем отдельный график для каждого изображения
for key, sizes in file_sizes.items():
    plt.figure(figsize=(10, 6))

    actual_quality_levels = list(range(0, 105, 5))
    plt.plot(actual_quality_levels, sizes, marker='o', color='blue', linewidth=2)

    plt.title(f'Зависимость размера сжатого файла от качества сжатия\n{key}')
    plt.xlabel('Коэффициент качества сжатия')
    plt.ylabel('Размер файла (байты)')
    plt.grid(True)

    # Сохраняем график в файл
    plt.savefig(f'{output_dir}/{key}_compression_graph.png')

print(f"Графики сохранены в директории {output_dir}")