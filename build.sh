#!/bin/bash

# Создаем директорию для бинарных файлов если её нет
mkdir -p cmd/client/bin

# Версия приложения
VERSION="1.0.1"

# Определяем переменную окружения для сервера клиента
SERVER_CLIENT_ADDRESS="87.228.37.67:8080"

# Функция для сборки
build() {
    local os=$1
    local arch=$2
    local extension=""
    
    if [ "$os" = "windows" ]; then
        extension=".exe"
    fi
    
    echo "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch go build \
    -o "cmd/client/bin/gophkeeper-client-${os}-${arch}${extension}" \
    -ldflags="-X 'main.VERSION=${VERSION}' \
              -X 'main.SERVER_CLIENT_ADDRESS=${SERVER_CLIENT_ADDRESS}'" \
    ./cmd/client/
    
    # Добавляем права на выполнение для Unix-систем
    if [ "$os" != "windows" ]; then
        chmod +x "cmd/client/bin/gophkeeper-client-${os}-${arch}${extension}"
    fi
}

# Функция для вывода меню
show_menu() {
    echo "Выберите платформу для сборки:"
    echo "1) Linux amd64"
    echo "2) Linux arm64"
    echo "3) macOS amd64"
    echo "4) macOS arm64"
    echo "5) Windows amd64"
    echo "6) Собрать для всех платформ"
    echo "0) Выход"
}

# Основной цикл
while true; do
    show_menu
    read -p "Введите номер (0-6): " choice
    
    case $choice in
        0)
            echo "Выход из программы"
            exit 0
            ;;
        1)
            build "linux" "amd64"
            ;;
        2)
            build "linux" "arm64"
            ;;
        3)
            build "darwin" "amd64"
            ;;
        4)
            build "darwin" "arm64"
            ;;
        5)
            build "windows" "amd64"
            ;;
        6)
            echo "Сборка для всех платформ..."
            build "linux" "amd64"
            build "linux" "arm64"
            build "darwin" "amd64"
            build "darwin" "arm64"
            build "windows" "amd64"
            ;;
        *)
            echo "Неверный выбор. Пожалуйста, выберите число от 0 до 6."
            ;;
    esac
    
    if [ $choice != 0 ]; then
        echo -e "\nСборка завершена. Бинарные файлы находятся в директории bin:"
        ls -l bin/
        echo -e "\nНажмите Enter для продолжения..."
        read
        clear
    fi
done 