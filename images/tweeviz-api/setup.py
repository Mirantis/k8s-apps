from setuptools import setup

setup(
    name='TweeViz-API',
    version='0.1.0',
    py_modules=['tweeviz_api'],
    include_package_data=True,
    install_requires=[
        'snakebite',
        'flask',
        'click',
        'cassandra-driver',
    ],
    entry_points='''
        [console_scripts]
        tweeviz-api=tweeviz_api:main
    ''',
)
