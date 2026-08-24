#include <stdio.h>

// Verify that a machine has ZMM20 support
int main() {
    // Zero out ZMM20 directly using EVEX encoding
    __asm__ volatile ("vpxord %zmm20, %zmm20, %zmm20");
    printf("ZMM20 written successfully!\n");
    return 0;
}
