#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <stdio.h>

int main() {
    struct stat stat;
    printf("macro sizeof(stat) %zu end\n", sizeof(stat));

    printf("macro sizeof(stat.st_dev) %zu end\n", sizeof(stat.st_dev));
    printf("macro sizeof(stat.st_dev) %zu end\n", sizeof(stat.st_dev));      
    printf("macro sizeof(stat.st_ino) %zu end\n", sizeof(stat.st_ino));      
    printf("macro sizeof(stat.st_mode) %zu end\n", sizeof(stat.st_mode));     
    printf("macro sizeof(stat.st_nlink) %zu end\n", sizeof(stat.st_nlink));    
    printf("macro sizeof(stat.st_uid) %zu end\n", sizeof(stat.st_uid));      
    printf("macro sizeof(stat.st_gid) %zu end\n", sizeof(stat.st_gid));      
    printf("macro sizeof(stat.st_rdev) %zu end\n", sizeof(stat.st_rdev));     
    printf("macro sizeof(stat.st_size) %zu end\n", sizeof(stat.st_size));     
    printf("macro sizeof(stat.st_blksize) %zu end\n", sizeof(stat.st_blksize));  
    printf("macro sizeof(stat.st_blocks) %zu end\n", sizeof(stat.st_blocks));   
    printf("macro sizeof(stat.st_atim) %zu end\n", sizeof(stat.st_atim));  
    printf("macro sizeof(stat.st_mtim) %zu end\n", sizeof(stat.st_mtim));  
    printf("macro sizeof(stat.st_ctim) %zu end\n", sizeof(stat.st_ctim));  

    return 0;
}
