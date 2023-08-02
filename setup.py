"""Arrrgs setup script"""
import setuptools

with open(".version", "r", encoding="utf-8") as fh:
    VERSION = ''.join(fh.read().split())

with open("README.md", "r", encoding="utf-8") as fh:
    README = fh.read()

setuptools.setup(
    name="hapm",
    version=VERSION,
    author="Mikhael Khrustik",
    description="The library for easily writing feature-rich Python scripts",
    long_description=README,
    long_description_content_type="text/markdown",
    packages=[
        'libhapm',
        'libhapm.cli',
        'libhapm.github',
        'libhapm.integration',
        'libhapm.package',
        'libhapm.manager',
        'libhapm.plugin',
        'libhapm.manifest',
    ],
    scripts=['scripts/hapm'],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License"
    ],
    install_requires=[
        'ruamel.yaml==0.17.21',
        'arrrgs==3.0.0',
        'PyGithub==1.59.0',
        'gitpython==3.1.30',
        'requests==2.31.0',
        'urllib3==1.26.6'
    ],
    python_requires='>=3.7',
    package_dir={'': '.'},
)
