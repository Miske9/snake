import os
import sys
import time
import random
import msvcrt

# Postavke igre
WIDTH = 20
HEIGHT = 10

def clear_screen():
    """Čisti terminal."""
    os.system('cls' if os.name == 'nt' else 'clear')

def create_board():
    """Kreiraj praznu ploču."""
    return [[" "] * WIDTH for _ in range(HEIGHT)]

def print_board(board):
    """Ispisuje ploču na konzolu."""
    clear_screen()
    print("+" + "-" * WIDTH + "+")
    for row in board:
        print("|" + "".join(row) + "|")
    print("+" + "-" * WIDTH + "+")

def place_food(board, snake):
    """Postavi hranu na ploču na slučajnu poziciju."""
    while True:
        x = random.randint(0, WIDTH - 1)
        y = random.randint(0, HEIGHT - 1)
        if (y, x) not in snake:
            board[y][x] = '*'
            return (y, x)

def move_snake(snake, direction):
    """Pomjera zmiju u zadanom smjeru."""
    head_y, head_x = snake[0]
    if direction == 'UP':
        new_head = (head_y - 1, head_x)
    elif direction == 'DOWN':
        new_head = (head_y + 1, head_x)
    elif direction == 'LEFT':
        new_head = (head_y, head_x - 1)
    elif direction == 'RIGHT':
        new_head = (head_y, head_x + 1)
    else:
        new_head = (head_y, head_x)
    return [new_head] + snake[:-1]

def check_collision(snake):
    """Provjera sudara zmije sa zidovima ili samom sobom."""
    head_y, head_x = snake[0]
    if head_y < 0 or head_y >= HEIGHT or head_x < 0 or head_x >= WIDTH:
        return True
    if snake[0] in snake[1:]:
        return True
    return False

def game_loop():
    """Glavna petlja igre."""
    snake = [(HEIGHT // 2, WIDTH // 2)]  # Početna pozicija zmije
    direction = 'RIGHT'
    board = create_board()
    food = place_food(board, snake)
    score = 0

    while True:
        # Prijem unosa korisnika
        if msvcrt.kbhit():
            key = msvcrt.getch().decode('utf-8').lower()
            if key == 'w' and direction != 'DOWN':
                direction = 'UP'
            elif key == 's' and direction != 'UP':
                direction = 'DOWN'
            elif key == 'a' and direction != 'RIGHT':
                direction = 'LEFT'
            elif key == 'd' and direction != 'LEFT':
                direction = 'RIGHT'

        # Pomiče zmiju
        new_snake = move_snake(snake, direction)

        # Provjera sudara
        if check_collision(new_snake):
            print("Game Over! Your score: ", score)
            print("Press 'r' to restart or 'q' to quit.")
            while True:
                if msvcrt.kbhit():
                    key = msvcrt.getch().decode('utf-8').lower()
                    if key == 'r':
                        game_loop()  # Pokreće igru ponovno
                        return
                    elif key == 'q':
                        print("Exiting the game...")
                        sys.exit()  # Izaći iz igre

        # Provjera pojela li je zmija hranu
        if new_snake[0] == food:
            new_snake.append(snake[-1])  # Produžava zmiju
            score += 1
            food = place_food(board, new_snake)  # Postavi novu hranu

        snake = new_snake

        # Ažuriraj ploču
        board = create_board()
        for y, x in snake:
            board[y][x] = 'o'
        board[food[0]][food[1]] = '*'

        print_board(board)
        print(f"Score: {score}")

        time.sleep(0.2)  # Pauza za brzinu igre

if __name__ == "__main__":
    game_loop()
