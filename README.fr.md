# LidReSleep

[English](./README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | **Français** | [Deutsch](./README.de.md) | [Español](./README.es.md) | [Русский](./README.ru.md) | [Português](./README.pt.md) | [Italiano](./README.it.md)

Un petit outil Windows en arrière-plan qui empêche votre ordinateur portable de chauffer pendant la nuit après la fermeture de l'écran.

De nombreux ordinateurs portables modernes (Modern Standby) ne se mettent pas vraiment en veille lorsque l'écran est fermé — l'écran s'éteint mais le système reste connecté à faible consommation, et peut facilement être **réveillé et maintenu éveillé** par des requêtes réseau, des tâches en arrière-plan, etc., provoquant chaleur et perte de batterie toute la nuit.

Ce que fait LidReSleep est simple : **veille à la fermeture de l'écran ; s'il est réveillé de façon inattendue alors que l'écran est encore fermé, il se rendort automatiquement jusqu'à ce que vous ouvriez l'écran.**

- Application portable en un seul fichier, sans installation
- 10 langues d'interface (简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano), détection automatique de la langue système
- Windows 10/11 (x64 / ARM64 / x86), choisissez la version correspondant à votre processeur

![Capture d'écran de LidReSleep](screenshot.jpg)

## Démarrage rapide

1. Téléchargez la version pour votre processeur, puis double-cliquez dessus (ex. `LidReSleep-amd64.exe`) pour ouvrir le panneau.
2. Cliquez sur **「▶ Démarrer la garde」** ; le statut passe à `● Garde active`.
3. Fermez l'écran et c'est tout.

Ensuite : fermer l'écran → veille immédiate ; réveillé avec l'écran encore fermé → veille à nouveau après ~3 secondes (par défaut) ; ouvrir l'écran → annulation et reprise normale.

## Guide de l'interface

### Statut
- `● Arrêté / ● Garde active`, indique si la garde est en cours.
- Bouton `▶ Démarrer la garde / ■ Arrêter la garde`, le texte bascule selon l'état.

### Paramètres

| Paramètre | Description | Défaut |
|---|---|---|
| Délai de sommeil (ms) | attendre ce délai avant de se rendormir après un réveil | `3000` |
| ☑ Lancer au démarrage | lancer automatiquement à la connexion Windows (niveau système, registre) | Non |
| ☑ Garde auto après connexion | démarrer la garde et réduire dans la zone de notification au lancement | Non |
| ☑ Réduire dans la zone de notification | masquer dans la zone de notification lors de la réduction | Oui |
| ☑ Fermer dans la zone de notification | masquer dans la zone de notification au lieu de quitter à la fermeture | Oui |

- Les modifications sont **enregistrées automatiquement**, aucune sauvegarde manuelle nécessaire.
- Le lancement au démarrage s'applique immédiatement (registre) ; les autres paramètres sont stockés dans `config.json` à côté de l'exe.

### Menus
- **Fichier** : Quitter
- **Outils** : Tester le sommeil (veille une fois pour vérifier)
- **Langue** : 10 langues (la coche indique la langue actuelle ; s'applique immédiatement)
- **Aide** : Rechercher des mises à jour (vérifie GitHub pour une nouvelle version), Page du projet, À propos (fonctionnalités et explication de Modern Standby)

### Journal
Chaque ligne comporte un horodatage et un niveau, affichant en temps réel les événements d'ouverture/fermeture, de réveil et de re-sommeil, avec défilement automatique et limite de 200 Ko.

| Niveau | Signification |
|---|---|
| `INFO` | informations générales |
| `EVENT` | événements système (écran/veille/réveil) |
| `ACTION` | actions du programme (planification/annulation/veille) |
| `ERROR` | erreur (avec raison) |

### Zone de notification
- Masquage dans la zone de notification lors de la réduction/fermeture (si activé).
- Clic gauche sur l'icône : restaurer la fenêtre principale ; clic droit : Afficher la fenêtre principale / Quitter.

## FAQ

**Pourquoi mon ordinateur portable chauffe-t-il encore après avoir fermé l'écran ?**
C'est probablement Modern Standby qui le réveille. Utilisez Outils → Tester le sommeil pour vérifier ; après avoir fermé l'écran pendant la garde, vérifiez le journal pour `ACTION Exécution du sommeil`.

**Qu'est-ce que Modern Standby ?**
Voir l'explication simple dans Aide → À propos.

**Comment quitter complètement ?**
Clic droit sur l'icône de la zone de notification → Quitter (si « Fermer dans la zone de notification » est activé, ✕ ne fait que réduire).

**Le lancement au démarrage ne fonctionne pas ?**
Cela dépend du compte utilisateur actuel ; les échecs (par exemple depuis un emplacement restreint) sont enregistrés avec une raison.

## Utilisateurs avancés : réduire les réveils (optionnel)

L'outil gère « se rendormir après un réveil ». Pour réduire les réveils eux-mêmes, exécutez dans un PowerShell élevé :

```powershell
# désactiver les minuteurs de réveil
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# inspecter les sources de réveil
powercfg /waketimers
powercfg /devicequery wake_armed
# restaurer les valeurs par défaut
powercfg /restoredefaultschemes
```

---

## Téléchargement

Obtenez la dernière version depuis GitHub Releases : [LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| Fichier | CPU | Téléchargement |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64 bits (la plupart des PC) | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64 (p. ex. Surface Pro X, PC Snapdragon) | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | x86 32 bits | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> Page du projet : https://github.com/dootn/LidReSleep
