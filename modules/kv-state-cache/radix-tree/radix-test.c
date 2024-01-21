#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include "radix.h"

typedef struct test_data {
    int k_index;
    int v_index;
} test_data;

// token list to insert
int *to_insert[] = {
    (int[]){103, 343, 123, 454, 0},
    (int[]){102, 343, 4564, 546, 342, 0},
    (int[]){103, 343, 4564, 546, 342, 0},
    (int[]){435, 7645, 4564, 546, 0},
    (int[]){435, 7645, 4564, 232, 454, 943, 0},
    (int[]){435, 343, 454, 123, 4533, 0},
    (int[]){435, 7645, 4564, 232, 454, 943, 0},
    NULL,
};

// token list to delete
int *to_remove[] = {
    (int[]){435, 7645, 4564, 546, 0},
    NULL,
};

// token list to find
int *to_find[] = {
    (int[]){435, 343, 454, 123, 4533, 0},
    (int[]){103, 343, 0},
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

size_t int_array_len(int *s) {
    size_t len = 0;
    while (s[len] != 0) len++;
    return len;
}

unsigned long insert_data_to_rax(rax *t, int *toadd[]) {
    unsigned long failed_insertions = 0;
    for (int i = 0; toadd[i] != nullptr; i++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = i;
        td->v_index = i;
        int retval = raxInsert(t, (unsigned char *)toadd[i], int_array_len(toadd[i]), (void *)td, NULL);
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

unsigned long insert_data_to_rax_and_return_node(rax *t, int *toadd[]) {
    unsigned len = 0;
    unsigned long failed_insertions = 0;
    for (int i = 0; toadd[i] != NULL; i++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = i;
        td->v_index = i;
        raxNode *node = raxInsertAndReturnDataNode(t, (unsigned char *)toadd[i], int_array_len(toadd[i]), (void *)td, NULL);
        if (node != NULL) {
            printf("Added token list: ");
            print_uint_array_as_string(toadd[i]);
            printf(", data: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
            struct test_data *new_td = (struct test_data *)malloc(sizeof(struct test_data));
            new_td->k_index = i+1;
            new_td->v_index = i+1;
            raxSetData(node,new_td);
        }
    }
    return failed_insertions;
}

// remove data from rax
void remove_data_from_rax(rax *t, int *toremove[]) {
    for (int i = 0; toremove[i] != NULL; i++) {
        if (raxRemove(t, (unsigned char *)toremove[i], int_array_len(toremove[i]), NULL)) {
            printf("raxRemove success, deleted token list: ");
        } else {
            printf("raxRemove failed for token list: ");
        }
        print_uint_array_as_string(toremove[i]);
        printf("\n");
    }
}

// find data in rax
void find_data_in_rax(rax *t, int *tofind[]) {
    for (int i = 0; tofind[i] != NULL; i++) {
        void *data = raxFind(t, (unsigned char *)tofind[i], int_array_len(tofind[i]));
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
    //insert_data_to_rax(t, to_insert);

    insert_data_to_rax_and_return_node(t, to_insert);

    // remove token list
    remove_data_from_rax(t, to_remove);

    // query token list
    find_data_in_rax(t, to_insert);

    raxFree(t);

    return 0;
}
