import vineyard
import numpy as np
import torch
import time


# create multiprocess here
import threading
from concurrent.futures import ProcessPoolExecutor, ThreadPoolExecutor
import random

def put(client, data):
    print("putting data in the thread:", threading.current_thread())
    random_name = random.randint(0, 10000)
    try:
        object_id = client.put(data, name="test" + str(random_name), persist=True)
        print("Successfully put the object_id:", object_id)
        client.release_object(object_id)
    except Exception as e:
        print("Failed to put the object", e)
        pass

def get(client):
    random_name = random.randint(0, 1000)
    try:
        client.get(name = "test" + str(random_name))
        print("Successfully get the object")
    except Exception as e:
        pass

def gc(interval=10):
    def gc_loop():
        client = vineyard.connect('/tmp/vineyard-local14.sock').fork()
        while True:
            # 获取所有对象名称
            names = client.list_names(pattern="*", limit=1000000)
            all_names = []
            object_ids = []
            for name, object_id in names.items():
                all_names.append(name)
                object_ids.append(object_id)

            if not object_ids:
                time.sleep(interval)  # 如果没有对象，休眠 interval 秒
                continue

            metas = client.get_metas(object_ids)
            blobs_to_delete = []

            # 遍历元数据获取 Blob
            for index, meta in enumerate(metas):
                blobs = vineyard.core.client._traverse_blobs(meta)
                for k, blob in blobs.items():
                    if not client.is_spilled(blob.id):
                        break
                    blobs_to_delete.append(all_names[index])

            for name in blobs_to_delete:
                client.delete_name(name=name)

            time.sleep(interval)  # 每次循环结束后休眠

    # 启动 GC 线程
    gc_thread = threading.Thread(target=gc_loop, daemon=True)
    gc_thread.start()

def put_thread():
    def put_loop():
        client = vineyard.connect('/tmp/vineyard-local14.sock').fork()
        data = torch.rand(100, 100, 100)
        while True:
            with ThreadPoolExecutor(max_workers=1) as executor:
                futures = []
                print("start the gc thread")
                for i in range(100):
                    futures.append(executor.submit(put, client, data))

                for future in futures:
                    future.result()
    put_thread = threading.Thread(target=put_loop, daemon=True)
    put_thread.start()

def get_thread():
    def get_loop():
        client = vineyard.connect('/tmp/vineyard-local14.sock').fork()
        while True:
            with ThreadPoolExecutor(max_workers=1) as executor:
                futures = []
                for i in range(100):
                    futures.append(executor.submit(get, client))

                for future in futures:
                    future.result()
    get_thread = threading.Thread(target=get_loop, daemon=True)
    get_thread.start()

put_thread()
#get_thread()
#gc(10)
import time
time.sleep(1000)
print("end")