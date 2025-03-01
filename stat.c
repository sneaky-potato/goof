#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <stdio.h>

int main() {
    printf("macro sizeof(stat) %zu end\n", sizeof(struct stat));

    st_dev;      
    st_ino;      
    st_mode;     
    st_nlink;    
    st_uid;      
    st_gid;      
    st_rdev;     
    st_size;     
    st_blksize;  
    st_blocks;   
    st_atim;  
    st_mtim;  
    st_ctim;  

    return 0;
}
