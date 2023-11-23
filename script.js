document.addEventListener('DOMContentLoaded', function () {
    const booksList = document.getElementById('books-list');
    const addBookForm = document.getElementById('add-book-form');

    function fetchBooks() {
        fetch('/api/books')
            .then(response => response.json())
            .then(data => {
                // Обновление списка книг на странице
                data.forEach(book => {
                    const bookItem = document.createElement('div');
                    bookItem.textContent = `${book.title} by ${book.author}`;
                    booksList.appendChild(bookItem);
                });
            })
            .catch(error => console.error('Error fetching books:', error));
    }

    function addBook(book) {
        fetch('/api/books', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(book),
        })
        .then(response => {
            if (response.ok) {
                // Обновить список книг после успешного добавления
                fetchBooks();
            } else {
                console.error('Error adding book:', response.statusText);
            }
        })
        .catch(error => console.error('Error adding book:', error));
    }

    // Вызов функции для получения списка книг при загрузке страницы
    fetchBooks();

    addBookForm.addEventListener('submit', function (event) {
        event.preventDefault();
        const title = document.getElementById('title').value;
        const author = document.getElementById('author').value;
        const genre = document.getElementById('genre').value;
        const description = document.getElementById('description').value;

        const newBook = {
            title: title,
            author: author,
            genre: genre,
            description: description,
        };

        addBook(newBook);
    });
});
