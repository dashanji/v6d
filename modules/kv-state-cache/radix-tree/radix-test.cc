#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include "radix.h"
//extern "C" {
//#include "../lz4/lib/lz4.h"
//}

struct test_data {
    int k_index;
    int v_index;
};

// token list to insert
int *to_insert[] = {
    /*(int[]){1, 1, 1, 0},
    (int[]){1, 1, 2, 0},
    (int[]){1, 1, 3, 0},
    (int[]){1, 2, 3, 4, 0},
    (int[]){2, 1, 1, 0},
    (int[]){1, 2, 0},
    (int[]){1, 2, 1, 0},
    (int[]){1, 2, 3, 0},
    (int[]){1, 2, 2, 0},
    (int[]){1, 2, 0},
    (int[]){2, 1, 0},
    (int[]){2, 2, 0},
    (int[]){2, 3, 10, 0},
    (int[]){2, 0},
    (int[]){1, 0},
    (int[]){1, 1, 0},
    (int[]){1, 2, 9, 0},
    (int[]){2, 2, 9, 0},
    (int[]){2, 3, 9, 0},
    (int[]){1, 2, 3, 5, 0},
    (int[]){2, 3, 9, 10, 0},
    (int[]){1, 2, 3, 5, 6, 0},*/
    (int[]){1, 0},
    (int[]){1, 2, 0},
    (int[]){1, 2, 3, 0},
    (int[]){1, 2, 3, 4, 0},
    (int[]){1, 2, 3, 4, 5, 0},
    //(int[]){2, 0},
    //(int[]){2, 6, 0},
    //(int[]){2, 8, 7, 0},
    //(int[]){2, 6, 7, 8, 0},
    //(int[]){2, 6, 7, 8, 9, 0},
    //(int[]){1, 2, 3, 0},
    //(int[]){1, 2, 0},
    //(int[]){1, 0},
    /* delete test case1*/
    //(int[]){1, 1, 0},
    //(int[]){2, 4, 0},
    //(int[]){3, 1, 0},
    /* delete test case2*/
    //(int[]){1, 1, 0},
    //(int[]){3, 1, 0},
    /* delete test case3*/
    //(int[]){3, 1, 2, 3, 5, 0},
    //(int[]){1, 1, 0},
    //(int[]){3, 1, 2, 4, 0},
    /*(int[]){1, 1, 1, 1, 3, 0},
    (int[]){1, 1, 2, 4, 6, 0},*/
    //(int[]){103, 343, 123, 454, 0},
    //(int[]){102, 343, 4564, 546, 342, 0},
    //(int[]){103, 343, 4564, 546, 342, 0},
    //(int[]){435, 7645, 4564, 546, 0},
    //(int[]){435, 7645, 4564, 232, 454, 943, 0},
    //(int[]){435, 343, 454, 123, 4533, 0},
    //(int[]){435, 7645, 4564, 232, 454, 943, 0},
    NULL,
};

// token list to delete
int *to_remove[] = {
    //(int[]){2, 3, 10, 0},
    //(int[]){2, 3, 9, 10, 0},
    /* delete test case1*/
    //(int[]){2,4, 0},
    /* delete test case2*/
    //(int[]){3, 1, 0},
    /* delete test case3*/
    //(int[]){3, 1, 2, 4, 0},
    //(int[]){2,1,1,0},
    //(int[]){1, 2, 9, 0},
    //(int[]){2, 2, 9, 0},
    //(int[]){2, 3, 9, 0},
    (int[]){1, 2, 3, 4, 0},
    NULL,
};

// token list to find
int *to_find[] = {
    /*(int[]){435, 343, 454, 123, 4533, 0},
    (int[]){103, 343, 0},*/
    (int[]){1, 2, 5, 0},
    (int[]){1, 2, 3, 5, 0},
    (int[]){1, 2, 3, 4, 6, 0},
    NULL,
};

void print_uint_array_as_string(int array[]) {
    printf("[");
    for (int i = 0; array[i] != 0; i++) {
        printf("%u", array[i]);
        if (array[i + 1] != 0) printf(", ");
    }
    printf("]");
}

void print_uint_array_with_len_as_string(int array[], int len) {
    printf("[");
    for (int i = 0; i<len; i++) {
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

unsigned long insert_data_to_rax(rax *t) {
    unsigned long failed_insertions = 0;
    for (int i = 0; to_insert[i] != NULL; i++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = i;
        td->v_index = i;
        int retval = raxInsert(t, (int *)to_insert[i], int_array_len(to_insert[i]), (void *)td, NULL);
        if (retval == 0) {
            if (errno == 0) {
                printf("Overwritten token list: ");
                print_uint_array_as_string(to_insert[i]);
                printf(", data: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
            } else {
                printf("Failed to insert for OOM:\n");
                print_uint_array_as_string(to_insert[i]);
            }
        } else {
            printf("Added token list: ");
            print_uint_array_as_string(to_insert[i]);
            printf(", data: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
        }
    }
    return failed_insertions;
}

/*
unsigned long insert_data_to_rax_and_return_node(rax *t) {
    unsigned len = 0;
    unsigned long failed_insertions = 0;
    for (int i = 0; to_insert[i] != NULL; i++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = i;
        td->v_index = i;
        raxNode *node = raxInsertAndReturnDataNode(t, (int *)to_insert[i], int_array_len(to_insert[i]), (void *)td, NULL);
        if (node != NULL) {
            printf("node is not null\n");
            struct test_data *new_td = (struct test_data *)malloc(sizeof(struct test_data));
            new_td->k_index = i+1;
            new_td->v_index = i+1;
            raxSetData(node,new_td);
            printf("Added token list: ");
            print_uint_array_as_string(to_insert[i]);
            printf(", new data: {k_index: %d, v_index: %d}\n", new_td->k_index, new_td->v_index);
        }
    }
    return failed_insertions;
}
*/

// remove data from rax
void remove_data_from_rax(rax *t) {
    for (int i = 0; to_remove[i] != NULL; i++) {
        if (raxRemove(t, (int *)to_remove[i], int_array_len(to_remove[i]), NULL)) {
            printf("raxRemove success, deleted token list: ");
        } else {
            printf("raxRemove failed for token list: ");
        }
        print_uint_array_as_string(to_insert[i]);
        printf("\n");
        raxShow(t);
        printf("t->numele: %d, t->numnodes: %d\n", t->numele, t->numnodes);
    }
}

// find data in rax
void find_data_in_rax(rax *t) {
    for (int i = 0; to_find[i] != NULL; i++) {
        void *data = raxFind(t, (int *)to_find[i], int_array_len(to_find[i]));
        if (data == raxNotFound) {
            printf("Token list ");
            print_uint_array_as_string(to_find[i]);
            printf(" is not found\n");
        } else {
            test_data *td = (struct test_data *)data;
            printf("Token list ");
            print_uint_array_as_string(to_find[i]);
            printf(" found, data is: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
        }
    }
}

/* Test the random walk function. */
int randomWalkTest(void) {
    rax *t = raxNew();
    long numele;
    for (numele = 0; to_insert[numele] != NULL; numele++) {
        raxInsert(t,(int*)to_insert[numele],
                    int_array_len(to_insert[numele]),(void*)numele,NULL);
    }

    raxIterator iter;
    raxStart(&iter,t);
    raxSeek(&iter,"^",NULL,0);
    int maxloops = 100000;
    while(raxRandomWalk(&iter,0) && maxloops--) {
        int nulls = 0;
        for (long i = 0; i < numele; i++) {
            if (to_insert[i] == NULL) {
                nulls++;
                continue;
            }
            if (int_array_len(to_insert[i]) == iter.key_len &&
                memcmp(to_insert[i],iter.key,iter.key_len) == 0)
            {
                to_insert[i] = NULL;
                nulls++;
            }
        }
        if (nulls == numele) break;
    }
    if (maxloops == 0) {
        printf("randomWalkTest() is unable to report all the elements "
               "after 100k iterations!\n");
        return 1;
    }
    raxStop(&iter);
    raxFree(t);
    return 0;
}

/* Iterator all tree */
int iteratorTest(){
    printf("++++++++++++++++ start iteratorTest ++++++++++++++++\n");
    rax *t = raxNew();
    long numele;
    for (numele = 0; to_insert[numele] != NULL; numele++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = numele;
        td->v_index = numele;
        raxInsert(t,(int*)to_insert[numele],
                    int_array_len(to_insert[numele]),(void*)td,NULL);
    }

    raxIterator iter;
    raxStart(&iter,t);
    raxSeek(&iter,"^",NULL,0);
    while(raxNext(&iter)) {
        test_data *td = (test_data *)(iter.data);
        printf("Token list ");
        print_uint_array_with_len_as_string(iter.key, iter.key_len);
        printf(" found, data is: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
    }
    raxStop(&iter);
    raxFree(t);
    return 0;
}

void raxSerializeTest() {
    printf("++++++++++++++++ start raxSerialize Test ++++++++++++++++\n");
    rax *t = raxNew();
    long numele;
    for (numele = 0; to_insert[numele] != NULL; numele++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = numele;
        td->v_index = numele;
        raxInsert(t,(int*)to_insert[numele],
                    int_array_len(to_insert[numele]),(void*)td,NULL);
    }
    size_t size = 0;
    std::vector<std::vector<int> > tokenList;
    std::vector<void*> dataList;
    std::vector<std::vector<int> > subtreeList;
    std::vector<void*> subtreeNodeList;

    printf("++++++++++++++++ original tree ++++++++++++++++\n");
    raxShow(t);
    printf("++++++++++++++++ start set subtree ++++++++++++++++\n");
    raxNode *parent = NULL;
    for (numele = 0; to_insert[numele] != NULL; numele++) {
        if (numele % 2 == 0) {
            continue;
        }
        print_uint_array_as_string(to_insert[numele]);
        raxNode *node = NULL;
        raxFindNodeWithParent(t, (int *)to_insert[numele], int_array_len(to_insert[numele]), (void **)&node, (void **)&parent);
        //raxFindNode(t, (int *)to_insert[numele], int_array_len(to_insert[numele]));
        node->issubtree = 1;
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = numele+100;
        td->v_index = numele+100;
        raxGetData(node);
        printf("in test function: parent: %p", parent);
        raxNode *newNode = raxReallocForTreeCustomData(node, (void **)parent);
        printf("####newNode is: %p\n", &newNode);
        raxSetCustomData(newNode, td);
    }
    raxShow(t);
    printf("++++++++++++++++ insert some new nodes ++++++++++++++++\n");
    for (numele = 0; to_find[numele] != NULL; numele++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = numele-100;
        td->v_index = numele-100;
        raxInsert(t,(int*)to_find[numele],
                    int_array_len(to_find[numele]),(void*)td,NULL);
    }

    raxShow(t);
    printf("++++++++++++++++ after raxSerialize ++++++++++++++++\n");
    raxSerialize(t, tokenList, dataList, &subtreeList, &subtreeNodeList);
    raxShow(t);
    printf("numele is: %d, numnodes: %d\n", t->numele, t->numnodes);
    
    printf("++++++++++++++++ start print subtree ++++++++++++++++\n");
    for (int i = 0; i < subtreeList.size(); i++) {
        raxNode *node = raxFindNode(t, subtreeList[i].data(), subtreeList[i].size());
        //node->issubtree = 0;
        printf("subtreeList[issubtree: %d]: ", node->issubtree);
        printf("node is %p\n", node);
        print_uint_array_with_len_as_string(subtreeList[i].data(), subtreeList[i].size());
        if (node->issubtree==1) {
            if (raxGetCustomData(node) == NULL) {
                printf(", customedata: null\n");
            } else {
                struct test_data *td = (struct test_data *)raxGetCustomData(node);
                printf(", customedata: {k_index: %d, v_index: %d}\n", td->k_index, td->v_index);
            }
        }
        printf("\n");
    }
    printf("++++++++++++++++ after set subtree ++++++++++++++++\n");
    raxShow(t);
    raxFree(t);
/*
  // We'll store some text into a variable pointed to by *src to be compressed later.
  const char* const src = serializedStr;
  // The compression function needs to know how many bytes exist.  Since we're using a string, we can use strlen() + 1 (for \0).
  const int src_size = size;
  // LZ4 provides a function that will tell you the maximum size of compressed output based on input data via LZ4_compressBound().
  const int max_dst_size = LZ4_compressBound(src_size);
  // We will use that size for our destination boundary when allocating space.
  char* compressed_data = (char*)malloc((size_t)max_dst_size);
  if (compressed_data == NULL)
    printf("Failed to allocate memory for *compressed_data.");
  // That's all the information and preparation LZ4 needs to compress *src into* compressed_data.
  // Invoke LZ4_compress_default now with our size values and pointers to our memory locations.
  // Save the return value for error checking.
  const int compressed_data_size = LZ4_compress_default(src, compressed_data, src_size, max_dst_size);
  // Check return_value to determine what happened.
  if (compressed_data_size <= 0)
    printf("A 0 or negative result from LZ4_compress_default() indicates a failure trying to compress the data. ");
  if (compressed_data_size > 0)
    printf("We successfully compressed some data! Ratio: %.2f\n",
        (float) compressed_data_size/src_size);
    printf("compressed_data_size is: %d, src_size is: %d\n", compressed_data_size, src_size);
  // Not only does a positive return_value mean success, the value returned == the number of bytes required.
  // You can use this to realloc() *compress_data to free up memory, if desired.  We'll do so just to demonstrate the concept.
  compressed_data = (char *)realloc(compressed_data, (size_t)compressed_data_size);
  if (compressed_data == NULL)
    printf("Failed to re-alloc memory for compressed_data.  Sad :(");


  // First, let's create a *new_src location of size src_size since we know that value.
  char* const regen_buffer = (char*)malloc(src_size);
  if (regen_buffer == NULL)
    printf("Failed to allocate memory for *regen_buffer.");
  // The LZ4_decompress_safe function needs to know where the compressed data is, how many bytes long it is,
  // where the regen_buffer memory location is, and how large regen_buffer (uncompressed) output will be.
  // Again, save the return_value.
  const int decompressed_size = LZ4_decompress_safe(compressed_data, regen_buffer, compressed_data_size, src_size);
  free(compressed_data); 
  if (decompressed_size < 0)
    printf("A negative result from LZ4_decompress_safe indicates a failure trying to decompress the data.  See exit code (echo $?) for value returned.");
  if (decompressed_size >= 0)
    printf("We successfully decompressed some data!\n");
  // Not only does a positive return value mean success,
  // value returned == number of bytes regenerated from compressed_data stream.
  if (decompressed_size != src_size)
    printf("Decompressed data is different from original! \n");


    rax *new_tree = raxDeserialize(regen_buffer);
    free(serializedStr);
    free(regen_buffer);
    find_data_in_rax(new_tree);
*/

}
/*
void raxSplitNodeTest(){
    printf("++++++++++++++++ start raxSplitNodeTest ++++++++++++++++\n");
    rax *t = raxNew();
    long numele;
    for (numele = 0; to_insert[numele] != NULL; numele++) {
        struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = numele;
        td->v_index = numele;
        raxInsert(t,(int*)to_insert[numele],
                    int_array_len(to_insert[numele]),(void*)td,NULL);
    }
    raxNode *node = (raxNode *)raxFind(t, (int *)to_insert[0], int_array_len(to_insert[0]));
    printf("node is: %p\n", node);
    struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
    td->k_index = 19;
    td->v_index = 19;
    
    raxNode** dataNodeList = (raxNode**) malloc(sizeof(raxNode*) * (t->numele));
    raxNode** current = dataNodeList;
    for (numele = 0; to_insert[numele] != NULL; numele++) {
                struct test_data *td = (struct test_data *)malloc(sizeof(struct test_data));
        td->k_index = numele;
        td->v_index = numele;
        printf("spliting..........\n");
        raxNode *new_node = raxNodeSplit(t,(int*)to_insert[numele], int_array_len(to_insert[numele]),(void*)td);
        printf("traversing.......... newnode:%d\n", new_node->data[0]);
        raxTraverseSubTree(new_node,&dataNodeList);
    }
    printf("traversing ererere..........\n");
    raxTraverse(t->head,&dataNodeList);
    raxShow(t);
    printf("numele is: %d, numnodes: %d\n", t->numele, t->numnodes);
    raxFree(t);
}
*/
int main() {
    /*
    rax *t = raxNew();
    if (t == NULL) return 1;

    // insert token list
    insert_data_to_rax(t);

    remove_data_from_rax(t);

    raxFree(t);
    */
   //raxSplitNodeTest();
 /*   
    rax *t = raxNew();
    if (t == NULL) return 1;

    // insert token list
    //insert_data_to_rax(t);
    insert_data_to_rax_and_return_node(t);
    printf("numele is: %d, numnodes: %d\n", t->numele, t->numnodes);
    // remove token list
    remove_data_from_rax(t);

    // query token list
    //find_data_in_rax(t);

    raxStack stack = raxFindWithStack(t, (int *)to_find[0], int_array_len(to_find[0]));
    printf("######stack items: %ld, max items: %ld######\n", stack.items, stack.maxitems);
    
    int items = stack.items;
    while (items > 0) {
        printf("items: %d\n", items);
        raxNode *node = raxStackPop(&stack);
        if (node == NULL) {
            printf("node is null\n");
        }
        printf("items: %ld, node: %p, numnodes: %d, nodeSize: %d, children: ", items, node, node->numnodes, node->size);
        for (int i = 0; i < node->size; i++) {
            printf(" %d", node->data[i]);
        }
        printf("\n");
        //test_data *td = (test_data *)raxGetData(node);
        //printf("node: %p, data is: {k_index: %d, v_index: %d}\n", node, td->k_index, td->v_index);
        items--;
    }
    printf("numele is: %d, numnodes: %d\n", t->numele, t->numnodes);

    //raxTraverse(t, raxShowCallback, NULL);
*/
    rax *t = raxNew();
    if (t == NULL) return 1;

    // insert token list
    insert_data_to_rax(t);

    remove_data_from_rax(t);

    raxFree(t);
    raxSerializeTest();
    //iteratorTest();
   //raxSerializeTest();
    //randomWalkTest();
    return 0;
}
