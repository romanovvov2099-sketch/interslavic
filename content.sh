#!/bin/sh

# Скрипт для конвертации проекта в txt файл
# Использование: ./project_to_txt.sh [директория_проекта] [выходной_файл]

# Устанавливаем параметры по умолчанию
PROJECT_DIR="${1:-.}"
OUTPUT_FILE="${2:-project_context.txt}"

# Расширения файлов для включения (можно изменить под ваш проект)
INCLUDE_EXTENSIONS="py js jsx ts tsx html css scss json yml yaml sql php java cpp c h go rb rs xml cfg ini conf dart"

# Директории для исключения
EXCLUDE_DIRS="node_modules venv env .git __pycache__ dist build .idea .vscode windows web macos linux ios android dart_tool build"

# Проверяем существование директории
if [ ! -d "$PROJECT_DIR" ]; then
    echo "Ошибка: Директория $PROJECT_DIR не найдена"
    exit 1
fi

echo "Конвертация проекта из $PROJECT_DIR в $OUTPUT_FILE..."

# Создаем временный файл
TEMP_FILE=$(mktemp)

# Формируем список исключаемых директорий для find
EXCLUDE_ARGS=""
for dir in $EXCLUDE_DIRS; do
    EXCLUDE_ARGS="$EXCLUDE_ARGS -not -path '*/$dir/*' 2>/dev/null"
done

# Формируем список включаемых расширений
EXT_ARGS=""
for ext in $INCLUDE_EXTENSIONS; do
    if [ -z "$EXT_ARGS" ]; then
        EXT_ARGS="-name '*.$ext'"
    else
        EXT_ARGS="$EXT_ARGS -o -name '*.$ext'"
    fi
done

# Заголовок файла
{
    echo "=========================================="
    echo "КОНТЕКСТ ПРОЕКТА"
    echo "Сгенерировано: $(date)"
    echo "Директория проекта: $(cd "$PROJECT_DIR" && pwd)"
    echo "=========================================="
    echo ""
} > "$TEMP_FILE"

# Добавляем структуру проекта
echo "СТРУКТУРА ПРОЕКТА:" >> "$TEMP_FILE"
echo "------------------------------------------" >> "$TEMP_FILE"

# Получаем структуру проекта (только файлы с нужными расширениями)
eval "find \"$PROJECT_DIR\" -type f $EXCLUDE_ARGS \( $EXT_ARGS \)" | while read -r file; do
    rel_path=$(echo "$file" | sed "s|^$PROJECT_DIR/||")
    echo "$rel_path" >> "$TEMP_FILE"
done

echo "" >> "$TEMP_FILE"
echo "==========================================" >> "$TEMP_FILE"
echo "СОДЕРЖИМОЕ ФАЙЛОВ:" >> "$TEMP_FILE"
echo "==========================================" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

# Счетчики для статистики
TOTAL_FILES=0
TOTAL_LINES=0

# Функция для добавления файла в вывод
process_file() {
    file="$1"
    rel_path=$(echo "$file" | sed "s|^$PROJECT_DIR/||")
    
    echo "Добавляю: $rel_path"
    
    # Добавляем разделитель и имя файла
    {
        echo ""
        echo "=========================================="
        echo "ФАЙЛ: $rel_path"
        echo "=========================================="
        echo ""
        
        # Пытаемся прочитать файл как текст
        if cat "$file" 2>/dev/null; then
            # Если успешно прочитали, добавляем перевод строки
            echo ""
        else
            echo "[Невозможно прочитать файл как текст]"
        fi
    } >> "$TEMP_FILE"
    
    # Подсчет статистики
    TOTAL_FILES=$((TOTAL_FILES + 1))
    if [ -r "$file" ]; then
        LINES=$(wc -l < "$file" 2>/dev/null || echo 0)
        TOTAL_LINES=$((TOTAL_LINES + LINES))
    fi
}

# Находим все файлы с нужными расширениями
echo "Поиск файлов..."

# Используем find с построенной командой и обрабатываем каждый файл
eval "find \"$PROJECT_DIR\" -type f $EXCLUDE_ARGS \( $EXT_ARGS \)" | while read -r file; do
    process_file "$file"
done

# Создаем финальный файл со статистикой в начале
FINAL_TEMP=$(mktemp)
{
    echo "=========================================="
    echo "КОНТЕКСТ ПРОЕКТА"
    echo "Сгенерировано: $(date)"
    echo "Директория проекта: $(cd "$PROJECT_DIR" && pwd)"
    echo "=========================================="
    echo ""
    echo "СТАТИСТИКА:"
    echo "- Обработано файлов: $TOTAL_FILES"
    echo "- Всего строк кода: $TOTAL_LINES"
    echo "- Размер файла: $(du -h "$TEMP_FILE" 2>/dev/null | cut -f1)"
    echo "=========================================="
    echo ""
    cat "$TEMP_FILE"
} > "$FINAL_TEMP"

# Перемещаем финальный файл
mv "$FINAL_TEMP" "$OUTPUT_FILE"

# Удаляем временный файл
rm -f "$TEMP_FILE"

echo ""
echo "Готово! Файл сохранен как: $OUTPUT_FILE"
echo "Статистика:"
echo "- Обработано файлов: $TOTAL_FILES"
echo "- Всего строк кода: $TOTAL_LINES"
echo "- Размер файла: $(du -h "$OUTPUT_FILE" 2>/dev/null | cut -f1)"