#include <stdio.h>

// Verify that a machine has ZMM20 support
int main() {
    // Check CPU features at runtime before invoking ZMM registers
    if (__builtin_cpu_supports("avx512f")) {
        __asm__ volatile ("vpxord %zmm20, %zmm20, %zmm20");
        printf("AVX-512 supported.\n");
    } else {
        printf("AVX-512 unsupported.\n");
    }
    return 0;
}
