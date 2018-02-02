from setuptools import setup

setup(
    name='TweePubFake',
    version='0.2.0',
    py_modules=['tweepub'],
    install_requires=[
        'click',
        'kafka-python==1.3.3',
    ],
    entry_points='''
        [console_scripts]
        tweepub=tweepub:main
    ''',
)
