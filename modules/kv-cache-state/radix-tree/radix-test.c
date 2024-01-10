#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include "radix.h"

typedef struct test_data {
    int k_index;
    int v_index;
} test_data;

// token list to insert
unsigned int *to_insert[] = {
    (unsigned int[]){103, 343, 123, 454, 0},
    (unsigned int[]){102, 343, 4564, 546, 342, 0},
    (unsigned int[]){103, 343, 4564, 546, 342, 0},
    (unsigned int[]){435, 7645, 4564, 546, 0},
    (unsigned int[]){435, 7645, 4564, 232, 454, 943, 0},
    (unsigned int[]){435, 343, 454, 123, 4533, 0},
    (unsigned int[]){435, 7645, 4564, 232, 454, 943, 0},
    NULL,
};

// token list to delete
unsigned int *to_remove[] = {
    (unsigned int[]){435, 7645, 4564, 546, 0},
    NULL,
};

// token list to find
unsigned int *to_find[] = {
    (unsigned int[]){435, 343, 454, 123, 4533, 0},
    (unsigned int[]){103, 343, 0},
    NULL,
};

void print_uint_array_as_string(unsigned int array[]) {
    printf("[");
    for (int i = 0; array[i] != 0; i++) {
        printf("%u", array[i]);
        if (array[i + 1] != 0) printf(", ");
    }
    printf("]");
}

size_t uint_array_len(unsigned int *s) {
    size_t len = 0;
    while (s[len] != 0) len++;
    return len;
}

unsigned long insert_data_to_rax(rax *t, unsigned int *toadd[]) {
    unsigned long failed_insertions = 0;
    for (int i = 0; toadd[i] != NULL; i++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = i;
        td->v_index = i;
        int retval = raxInsert(t, (unsigned char *)toadd[i], uint_array_len(toadd[i]), (void *)td, NULL);
        if (retval == 0) {
            if (errno == 0) {
                printf("Overwritten token list: ");
                print_uint_array_as_string(toadd[i]);
                printf(", data: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
            } else {
                printf("Failed to insert for OOM:\n");
                print_uint_array_as_string(toadd[i]);
            }
        } else {
            printf("Added token list: ");
            print_uint_array_as_string(toadd[i]);
            printf(", data: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
        }
    }
    return failed_insertions;
}

// remove data from rax
void remove_data_from_rax(rax *t, unsigned int *toremove[]) {
    for (int i = 0; toremove[i] != NULL; i++) {
        if (raxRemove(t, (unsigned char *)toremove[i], uint_array_len(toremove[i]), NULL)) {
            printf("raxRemove success, deleted token list: ");
        } else {
            printf("raxRemove failed for token list: ");
        }
        print_uint_array_as_string(toremove[i]);
        printf("\n");
    }
}

// find data in rax
void find_data_in_rax(rax *t, unsigned int *tofind[]) {
    for (int i = 0; tofind[i] != NULL; i++) {
        void *data = raxFind(t, (unsigned char *)tofind[i], uint_array_len(tofind[i]));
        if (data == raxNotFound) {
            printf("Token list ");
            print_uint_array_as_string(tofind[i]);
            printf(" is not found\n");
        } else {
            test_data *td = (test_data *)data;
            printf("Token list ");
            print_uint_array_as_string(tofind[i]);
            printf(" found, data is: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
        }
    }
}

int main() {
    rax *t = raxNew();
    if (t == NULL) return 1;

    // insert token list
    insert_data_to_rax(t, to_insert);

    // remove token list
    remove_data_from_rax(t, to_remove);

    // query token list
    find_data_in_rax(t, to_find);

    raxFree(t);
    return 0;
}

/*
Suppose the maximum size of a vineyard blob is a const MAX.(index < 1024)

Cache_builder will handle the mapping between vineyard blob and token kv_states.


Insert token list: [103, 343, 123, 454], kv_states_value is std::map<int, std::vector<double>, std::vector<double>>

1. builder = new Cache_builder(kv_states_value)

raxInsert(rax *rax, unsigned int *s, size_t len, void *data)

2. raxInsert(rax *rax, unsigned int *s, size_t len, builder);
{
    insert token list to rax, and the data is kv_states.

    builder mark the 

}

If the value < MAX, we can store it in one vineyard blob.
If the value > MAX, we can store it in multiple vineyard blobs.




Common case:

case 1:
    Insert a token list, and a vineyard object can store all kv states of the token list.
    [103, 343, 123, 454, 0]


Special case:

case 1:
    Insert a token list, but a vineyard object can't store all kv states of the token list.
    [103, 434, 343, 232, 2334, 343, 2323, ...]
case 2:

*/