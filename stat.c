#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <stdio.h>
#include <stddef.h>

int main() {
    struct stat stat;
    printf("macro sizeof(stat) %zu end\n", sizeof(stat));

    printf("macro stat.st_dev %zu + end\n", offsetof(struct stat, st_dev));      
    printf("macro stat.st_ino %zu + end\n", offsetof(struct stat, st_ino));      
    printf("macro stat.st_mode %zu + end\n", offsetof(struct stat, st_mode));     
    printf("macro stat.st_nlink %zu + end\n", offsetof(struct stat, st_nlink));    
    printf("macro stat.st_uid %zu + end\n", offsetof(struct stat, st_uid));      
    printf("macro stat.st_gid %zu + end\n", offsetof(struct stat, st_gid));      
    printf("macro stat.st_rdev %zu + end\n", offsetof(struct stat, st_rdev));     
    printf("macro stat.st_size %zu + end\n", offsetof(struct stat, st_size));     
    printf("macro stat.st_blksize %zu + end\n", offsetof(struct stat, st_blksize));  
    printf("macro stat.st_blocks %zu + end\n", offsetof(struct stat, st_blocks));   
    printf("macro stat.st_atim %zu + end\n", offsetof(struct stat, st_atim));  
    printf("macro stat.st_mtim %zu + end\n", offsetof(struct stat, st_mtim));  
    printf("macro stat.st_ctim %zu + end\n", offsetof(struct stat, st_ctim));  

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
