from setuptools import setup

setup(
    name='TweePub',
    version='0.2.0',
    py_modules=['tweepub'],
    install_requires=[
        'click',
        'tweepy',
        'kafka-python==1.3.3',
        'pyhdfs'
    ],
    entry_points='''
        [console_scripts]
        tweepub=tweepub:main
    ''',
)
