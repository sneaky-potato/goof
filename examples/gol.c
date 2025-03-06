#include <stdio.h>
#include <unistd.h>

#define ROWS 16
#define COLS 32

int front[ROWS][COLS] = {0};

void display() {
    for (int y = 0; y < ROWS; y++) {
        for (int x = 0; x < COLS; x++) {
            if (front[y][x]) {
                printf("#");
            } else {
                printf(".");
            }
        }
        printf("\n");
    }
}

int mod(int a, int b) {
    return (a % b + b) % b;
}

int count_nbors(int cx, int cy) {
    int nbors = 0;
    for (int dx = -1; dx <= 1; dx++) {
        for (int dy = -1; dy <= 1; dy++) {
            if (!(dx == 0 && dy == 0)) {
                int x = mod(cx + dx, COLS);
                int y = mod(cy + dy, ROWS);
                if ((front[y][x] == 1) || (front[y][x] == 3)) nbors += 1;
            }
        }
    }
    return nbors;
}

void next() {
    for (int y = 0; y < ROWS; y++) {
        for (int x = 0; x < COLS; x++) {
            int nbors = count_nbors(x, y);
            if (front[y][x] == 1 && (nbors < 2 || nbors > 3)) {
                front[y][x] = 3;
            } else if (front[y][x] == 0 && nbors == 3) {
                front[y][x] = 2;
            }
        }
    }

    for (int y = 0; y < ROWS; y++) {
        for (int x = 0; x < COLS; x++) {
            if (front[y][x] == 2) front[y][x] = 1;
            if (front[y][x] == 3) front[y][x] = 0;
        }
    }
}

int main() {
    // 010
    // 001
    // 111
    front[0][1] = 1;
    front[1][2] = 1;
    front[2][0] = 1;
    front[2][1] = 1;
    front[2][2] = 1;
    for (;;) {
        next();
        display();
        printf("\033[%dA\033[%dD", ROWS, COLS);
        usleep(100 * 1000);
    }
    return 0;
}
