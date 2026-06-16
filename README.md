# Nexora

Forum web moderne développé en Go

## Présentation

Nexora est une application web de type forum permettant aux utilisateurs de créer un compte, publier du contenu, partager des images et interagir avec les autres membres grâce à un système de commentaires et de réactions.

Ce projet a été réalisé dans le cadre de la première année de Bachelor Informatique afin de mettre en pratique les notions de développement web full-stack, de gestion de base de données et d'architecture logicielle.

---

## Fonctionnalités

### Gestion des utilisateurs

* Inscription
* Connexion
* Déconnexion
* Gestion des sessions

### Publications

* Création de posts
* Ajout d'une image
* Attribution d'une catégorie
* Affichage des publications

### Interactions

* Commentaires
* Likes
* Dislikes

### Recherche et filtrage

* Recherche de publications
* Filtre par catégorie

### Interface utilisateur

* Design moderne
* Affichage dynamique des données
* Mise à jour des publications sans rechargement complet de la page

---

## Technologies utilisées

### Backend

* Go (Golang)
* API REST
* Gestion des routes HTTP

### Frontend

* HTML
* CSS
* JavaScript

### Base de données

* SQLite

### Outils

* Git
* GitHub
* Docker

---

## Architecture du projet

```text
forum-project
│
├── cmd
│   └── server
│       └── main.go
│
├── internal
│   ├── database
│   │   └── database.go
│   │
│   └── models
│
├── static
│   ├── css
│   │   └── style.css
│   │
│   ├── js
│   │   └── main.js
│   │
│   └── uploads
│
├── templates
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   └── create_post.html
│
├── forum.db
│
└── README.md
```

---

## Base de données

L'application repose sur plusieurs tables principales :

### Users

Contient les informations des utilisateurs :

* id
* username
* email
* password

### Posts

Contient les publications :

* id
* user_id
* title
* content
* image_path

### Categories

Contient les catégories de publication :

* Général
* Développement
* Aide
* Projet
* Discussion
* Sport
* Lifestyle
* Musique
* Gaming
* Cinéma
* Voyage
* Études

### Comments

Contient les commentaires associés aux publications.

### Reactions

Contient les likes et dislikes des utilisateurs.

---

## Installation

### Prérequis

* Go 1.24 ou supérieur
* Git

### Cloner le projet

```bash
git clone https://github.com/VOTRE-USERNAME/forum-project.git

cd forum-project
```

### Installer les dépendances

```bash
go mod tidy
```

### Lancer l'application

```bash
go run cmd/server/main.go
```

Le serveur sera accessible à l'adresse :

```text
http://localhost:8080
```

---

## Utilisation

### Créer un compte

1. Cliquer sur "Inscription"
2. Renseigner :

   * nom d'utilisateur
   * email
   * mot de passe

### Se connecter

1. Cliquer sur "Connexion"
2. Entrer ses identifiants

### Créer une publication

1. Cliquer sur "Créer un post"
2. Ajouter :

   * un titre
   * un contenu
   * une catégorie
   * une image (optionnelle)

### Interagir

* Ajouter un commentaire
* Liker une publication
* Disliker une publication

---

## Captures d'écran

Ajouter ici :

### Accueil

![Accueil](captures/home.png)

### Création d'un post

![Création](captures/create-post.png)

### Publication avec image

![Post](captures/post-image.png)

---

## Difficultés rencontrées

Durant le développement du projet, plusieurs défis ont été rencontrés :

* Mise en place de l'authentification
* Gestion des sessions utilisateurs
* Upload et affichage des images
* Communication entre le backend Go et le frontend JavaScript
* Mise à jour dynamique des publications

Ces difficultés ont permis d'approfondir la compréhension du développement web full-stack.

---

## Améliorations possibles

* Profils utilisateurs
* Notifications
* Messagerie privée
* Administration des contenus
* Responsive mobile avancé
* Pagination des publications

---

## Auteur

Pauline Tumbarello

Bachelor Informatique – Première année

Projet réalisé dans le cadre de l'apprentissage du développement web full-stack avec Go.
