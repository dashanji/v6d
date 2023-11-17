from kfp import dsl

def PreProcess(data_multiplier: int, registry: str):
    op =  dsl.ContainerOp(
        name='Preprocess Data',
        image = f'{registry}/preprocess-data',
        container_kwargs={
            'image_pull_policy': "Always",
        },
        pvolumes={
            "/data": dsl.PipelineVolume(pvc="benchmark-data"),
        },
        command = ['python3', 'preprocess.py'],
        arguments = [f'--data_multiplier={data_multiplier}'],
    )
    return op

def Train(comp1, registry: str):
    op = dsl.ContainerOp(
        name='Train Data',
        image=f'{registry}/train-data',
        container_kwargs={
            'image_pull_policy': "Always",
        },
        pvolumes={
            "/data": comp1.pvolumes['/data'],
        },
        command = ['python3', 'train.py'],
    )
    return op

def Test(comp2, registry: str):
    op = dsl.ContainerOp(
        name='Test Data',
        image=f'{registry}/test-data',
        container_kwargs={
            'image_pull_policy': "Always",
        },
        pvolumes={
            "/data": comp2.pvolumes['/data'],
        },
        command = ['python3', 'test.py'],
    )
    return op

@dsl.pipeline(
   name='Machine Learning Pipeline',
   description='An example pipeline that trains and logs a regression model.'
)
def pipeline(data_multiplier: int, registry: str):
    comp1 = PreProcess(data_multiplier=data_multiplier, registry=registry)
    comp2 = Train(comp1, registry=registry)
    comp3 = Test(comp2, registry=registry)

if __name__ == '__main__':
    from kfp import compiler
    compiler.Compiler().compile(pipeline, __file__[:-3]+ '.yaml')
