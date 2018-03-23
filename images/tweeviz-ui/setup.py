from setuptools import setup

setup(
    name='TweeViz-UI',
    version='0.1.0',
    packages=['tweeviz_ui'],
    include_package_data=True,
    install_requires=[
        'flask',
    ],
    entry_points='''
        [console_scripts]
        tweeviz-ui=tweeviz_ui:main
    ''',
)
