from setuptools import setup

setup(
    name='TweePub',
    version='0.1.0',
    py_modules=['tweepub'],
    install_requires=[
        'click',
        'tweepy',
        'kafka-python',
    ],
    entry_points='''
        [console_scripts]
        tweepub=tweepub:main
    ''',
)
