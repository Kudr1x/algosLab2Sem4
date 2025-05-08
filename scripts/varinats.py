from PIL import Image
import numpy as np
import os


def convert_to_grayscale(input_image):
    """Преобразует изображение в оттенки серого"""
    gray_image = input_image.convert("L")
    return gray_image


def convert_to_bw_no_dithering(input_image):
    """Преобразует изображение в черно-белое без дизеринга"""
    # Используем метод convert("1") для получения чисто черно-белого изображения с порогом 127
    bw_image = input_image.convert("1")
    return bw_image


def bayer_matrix(n):

    """Создает матрицу Байера размером 2^n x 2^n"""
    if n == 0:
        return np.array([[0]])
    else:
        bm = bayer_matrix(n - 1)
        size = 2 ** (n - 1)
        return np.block([
            [4 * bm, 4 * bm + 2],
            [4 * bm + 3, 4 * bm + 1]
        ])


def convert_to_bw_with_dithering(input_image, bayer_level=2):
    """Преобразует изображение в черно-белое с дизерингом Байера"""
    # Сначала преобразуем в оттенки серого
    gray_image = convert_to_grayscale(input_image)
    gray_array = np.array(gray_image)

    # Создаем матрицу Байера выбранного размера
    matrix_size = 2 ** bayer_level
    threshold_map = bayer_matrix(bayer_level)

    # Нормализуем матрицу Байера до диапазона 0-255
    threshold_map = threshold_map * (256 / (matrix_size ** 2))

    # Применяем дизеринг
    height, width = gray_array.shape
    for y in range(height):
        for x in range(width):
            # Находим соответствующую ячейку в матрице Байера
            x_map = x % matrix_size
            y_map = y % matrix_size

            # Если значение пикселя больше порога - белый, иначе - черный
            if gray_array[y, x] > threshold_map[y_map, x_map]:
                gray_array[y, x] = 255
            else:
                gray_array[y, x] = 0

    # Создаем новое изображение из массива
    dithered_image = Image.fromarray(gray_array.astype(np.uint8))
    return dithered_image


def process_image(input_path, output_dir, bayer_level=2):
    """Обрабатывает одно изображение и сохраняет результаты"""
    try:
        # Создаем выходную директорию, если она не существует
        os.makedirs(output_dir, exist_ok=True)

        # Получаем имя файла без расширения
        filename = os.path.basename(input_path)
        name_without_ext = os.path.splitext(filename)[0]

        # Открываем исходное изображение
        image = Image.open(input_path)

        # Преобразуем в оттенки серого
        gray_image = convert_to_grayscale(image)
        gray_image.save(os.path.join(output_dir, f"{name_without_ext}_grayscale.png"))
        print(f"Изображение в оттенках серого сохранено: {name_without_ext}_grayscale.png")

        # Преобразуем в черно-белое без дизеринга
        bw_image = convert_to_bw_no_dithering(image)
        bw_image.save(os.path.join(output_dir, f"{name_without_ext}_bw_no_dithering.png"))
        print(f"Черно-белое изображение без дизеринга сохранено: {name_without_ext}_bw_no_dithering.png")

        # Преобразуем в черно-белое с дизерингом
        dithered_image = convert_to_bw_with_dithering(image, bayer_level)
        dithered_image.save(os.path.join(output_dir, f"{name_without_ext}_bw_dithering.png"))
        print(f"Черно-белое изображение с дизерингом сохранено: {name_without_ext}_bw_dithering.png")

        return True
    except Exception as e:
        print(f"Ошибка при обработке {input_path}: {e}")
        return False


def main():
    # Задаем параметры обработки
    # Список изображений для обработки
    image_paths = [
        "/home/kudrix/GolandProjects/AlgosLab2Sem4v4/varinats/forest.png",
    ]

    # Директория для сохранения результатов
    output_directory = "/home/kudrix/GolandProjects/AlgosLab2Sem4v4/varinats"

    # Уровень матрицы Байера для дизеринга (1-4)
    bayer_level = 2

    # Обрабатываем каждое изображение
    successful = 0
    for image_path in image_paths:
        if process_image(image_path, output_directory, bayer_level):
            successful += 1

    print(f"Обработка завершена. Успешно обработано {successful} из {len(image_paths)} изображений.")


if __name__ == "__main__":
    main()
