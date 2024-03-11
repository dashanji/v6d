import json

def generate_long_prompts():
    f = open('/opt/caoye/ShareGPT_Vicuna_unfiltered/ShareGPT_V3_unfiltered_cleaned_split.json', 'r')
    data = json.loads(f.read())

    human_values = []

    prompts_file = open('/opt/caoye/suri.cpp/examples/llama-with-vineyard/share-gpt-prompts.txt', 'w')
    prompts = []
    prompts_len = 0
    for entry in data:
        prefix_prompt = ""
        if len(entry.get("conversations", [])) < 6:
            continue
        for conversation in entry.get("conversations", []):
            if conversation.get("from") == "gpt":
                prefix_prompt += conversation.get("value").replace("\n", " ")
            if conversation.get("from") == "human":
                prefix_prompt += conversation.get("value").replace("\n", " ")
                #prefix_prompt.replace("\n", " ").replace("\n", " ")
                prefix_prompt.replace("\n", " ")
                prompts_file.write(prefix_prompt + "\n")
                prompts_len += 1
                if prompts_len > 1000:
                    exit()

def generate_short_prompts():
    f = open('/opt/caoye/ShareGPT_Vicuna_unfiltered/ShareGPT_V3_unfiltered_cleaned_split.json', 'r')
    data = json.loads(f.read())

    human_values = []

    prompts_file = open('/opt/caoye/suri.cpp/examples/llama-with-vineyard/share-gpt-short-prompts.txt', 'w')
    prompts_len = 0
    for entry in data:
        prefix_prompt = ""
        prompts = []
        if len(entry.get("conversations", [])) < 6:
            continue
        for conversation in entry.get("conversations", []):
            if conversation.get("from") == "gpt":
                prompts.append(conversation.get("value").replace("\n", " "))
            if conversation.get("from") == "human":
                prompts.append(conversation.get("value").replace("\n", " "))
        count = 0
        for prompt in prompts:
            count += len(prompt)
        if count < 1000:
            prefix_prompt = ""
            for prompt in prompts:
                prompts_file.write(prefix_prompt + prompt + "\n")
                prefix_prompt += prompt + " "
            prompts_len += 1
        if prompts_len > 1000:
            exit()

generate_short_prompts()
