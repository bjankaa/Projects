# AlexNet Image Classification

An image classification project exploring the AlexNet convolutional neural
network architecture on the CIFAR-10 dataset.

The project adapts an AlexNet-based model to the smaller 32×32 CIFAR-10
images and evaluates different training and optimization strategies.

## Overview

The goal of this project was to study the architecture and training process
of convolutional neural networks through a practical image classification task.

The original AlexNet architecture was designed for ImageNet, so the model
was adapted to work with the CIFAR-10 dataset and its 10 image classes.

## Features

- AlexNet-based convolutional neural network
- CIFAR-10 image classification
- Convolutional and fully connected layers
- Batch normalization
- Max pooling
- Dropout regularization
- Data augmentation
- Training and validation pipeline
- Confusion matrix evaluation
- Comparison of multiple optimizers

## Technologies

- Python
- Keras
- Jupyter Notebook
- CIFAR-10

## Model Training

The dataset is divided into training and validation sets and processed before
training.

Data augmentation is used to introduce variation into the training images,
while callbacks are used to adjust the learning rate based on validation
performance.

The model uses classification accuracy as an evaluation metric.

## Optimizer Comparison

Multiple optimizers were evaluated to compare their effect on classification
performance.

| Optimizer | Accuracy |
|-----------|---------:|
| Adam      | 69.89%   |
| SGD       | 70.88%   |
| Nadam     | 73.36%   |

Nadam achieved the best accuracy among the evaluated optimizers.

## Evaluation

Model performance was evaluated using:

- Training and validation accuracy
- Training and validation loss
- Confusion matrices
- Classification accuracy

The confusion matrix showed that some visually similar CIFAR-10 classes,
particularly cats and dogs, were more difficult for the model to distinguish.

## Repository Contents

- `AlexNet_eredeti.ipynb` — Jupyter Notebook containing the model and experiments
- `AlexNet.pptx` — presentation about the project
- `temalabor_Bacso_Janka.pdf` — detailed project report

## Key Takeaways

This project provided practical experience with convolutional neural networks,
image preprocessing, model training, regularization techniques, and optimizer
selection.

It also demonstrated how architectural and training parameters can significantly
affect the performance of a neural network.