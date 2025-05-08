import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import os

# Загружаем данные из CSV
df = pd.read_csv('/home/kudrix/GolandProjects/AlgosLab2Sem4v4/compression_results.csv')

# Получаем список уникальных изображений
images = df['Image'].unique()

# Создаем график
plt.figure(figsize=(12, 8))

# Устанавливаем логарифмический масштаб для оси Y
plt.yscale('log')

# Уникальные цвета и маркеры для каждого изображения
colors = ['blue', 'red', 'green', 'purple']
markers = ['o', 's', '^', 'D']

# Строим график для каждого изображения
for i, image in enumerate(images):
    data = df[df['Image'] == image]
    plt.plot(
        data['Quality'],
        data['CompressedSize']/1024, # Переводим в килобайты
        label=image,
        color=colors[i % len(colors)],
        marker=markers[i % len(markers)],
        markersize=6,
        linewidth=2
    )

# Добавляем заголовок и подписи осей
plt.title('Зависимость размера сжатого файла от качества сжатия ', fontsize=16)
plt.xlabel('Качество сжатия', fontsize=14)
plt.ylabel('Размер файла (КБ, log scale)', fontsize=14)
plt.grid(True, linestyle='--', alpha=0.7)
plt.legend(fontsize=12)

# Устанавливаем диапазон для оси X
plt.xlim(0, 100)

# Добавляем горизонтальные линии сетки для логарифмического масштаба
plt.grid(True, which="both", ls="-", alpha=0.2)
plt.grid(True, which="major", ls="-", alpha=0.5)

# Аннотации для экстремальных точек
for image in images:
    data = df[df['Image'] == image]
    min_size = data['CompressedSize'].min()
    max_size = data['CompressedSize'].max()

    min_q = data.loc[data['CompressedSize'] == min_size, 'Quality'].values[0]
    max_q = data.loc[data['CompressedSize'] == max_size, 'Quality'].values[0]

    # Для логарифмического масштаба смещение аннотаций нужно делать по-другому
    plt.annotate(f'{min_size/1024:.1f} КБ',
                 xy=(min_q, min_size/1024),
                 xytext=(min_q-10, min_size/1024*0.9),
                 arrowprops=dict(arrowstyle='->'))

    plt.annotate(f'{max_size/1024:.1f} КБ',
                 xy=(max_q, max_size/1024),
                 xytext=(max_q+10, max_size/1024*1.1),
                 arrowprops=dict(arrowstyle='->'))

# Делаем более информативную разметку для логарифмической оси
plt.minorticks_on()

# Сохраняем график
plt.tight_layout()
plt.savefig('compression_chart_log_scale.png', dpi=300, bbox_inches='tight')
plt.show()
