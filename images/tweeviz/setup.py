from setuptools import setup

setup(
    name='TweeViz',
    version='0.2.0',
    py_modules=['tweeviz'],
    include_package_data=True,
    install_requires=[
        'snakebite',
        'flask',
        'click',
        'cassandra-driver',
    ],
    entry_points='''
        [console_scripts]
        tweeviz=tweeviz:main
    ''',
)
