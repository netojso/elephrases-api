# Elephrases - Technical Requirements

## 1. Overview
A flashcards app focused on language learning, with support for text and media files, packaged in `.elephrases` format.

## 2. Core Features
- [x] Create and edit flashcards with text and optional images/audio.
- [x] Group flashcards into decks.
- [ ] Export decks as `.elephrases` files (a compressed format).
- [ ] Import `.elephrases` files to restore decks.
- [ ] Upload and download `.elephrases` files via HTTP endpoints.
- [x] Study mode with predefined settings.

### Study Mode
- [x] Allow users to enter a study session for any deck.
- [ ] Predefined settings:
  - Interval between cards (e.g., 3s, 5s, 10s)
  - Shuffle cards randomly
  - Repeat missed cards until answered correctly
  - Session length (e.g., 10, 20, all cards)
  - Option to show only front or both front/back
- [x] Track progress within a session (e.g., correct/incorrect count)
- [ ] End-of-session summary with stats

## 3. Backend (Golang + Gin)
- [x] REST API for creating, reading, updating, and deleting flashcards.
- [x] Multipart file upload support for media files.
- [ ] Service to package decks and media into a `.elephrases` file (ZIP).
- [ ] Service to extract and import `.elephrases` files.
- [x] Database migrations system (`database/migrations`).

## 4. Authentication
- [x] JWT-based user authentication
- [x] Register and login endpoints
- [x] Protected routes for creating and modifying flashcards
- [x] Token refresh mechanism (optional)

## 5. Frontend / Client
- [ ] Web or mobile client to manage flashcards.
- [ ] Interface to upload/download `.elephrases` files.
- [ ] Display media content within flashcards.
- [ ] Login and user session management
- [ ] Study mode UI with session control and feedback

## 6. File Format (.elephrases)
- [ ] ZIP archive containing:
  - `deck.json` (flashcard data)
  - `media/` folder with all images/audio used

## 7. Future Ideas
- [ ] Tag system for flashcards
- [ ] Spaced repetition algorithm
- [ ] User accounts and sync

---

*Last updated: 2025-04-25*
