document.addEventListener("DOMContentLoaded", () => {
    loadPosts();
});

async function loadPosts() {
    const postsList = document.getElementById("postsList");
    const searchInput = document.getElementById("searchInput");
    const categoryFilter = document.getElementById("categoryFilter");

    const postsCount = document.getElementById("postsCount");
    const usersCount = document.getElementById("usersCount");
    const commentsCount = document.getElementById("commentsCount");

    if (!postsList) {
        return;
    }

    try {
        const response = await fetch("/api/posts");
        const posts = await response.json();

        let totalComments = 0;

        posts.forEach(post => {
            totalComments += post.comments ? post.comments.length : 0;
        });

        postsCount.textContent = posts.length;
        usersCount.textContent = "-";
        commentsCount.textContent = totalComments;

        function displayPosts(filteredPosts) {
            postsList.innerHTML = "";

            if (filteredPosts.length === 0) {
                postsList.innerHTML = `
                    <article class="post-card">
                        <h3>Aucun post</h3>
                        <p>Aucun post trouvé.</p>
                    </article>
                `;
                return;
            }

            filteredPosts.forEach(post => {
                const article = document.createElement("article");
                article.className = "post-card";
                article.dataset.title = post.title.toLowerCase();
                article.dataset.category = post.category;

                let imageHTML = "";

                if (post.image && post.image !== "") {
                    imageHTML = `
                        <img class="post-image" src="${post.image}" alt="Image du post">
                    `;
                }

                let commentsHTML = "";

                if (post.comments && post.comments.length > 0) {
                    post.comments.forEach(comment => {
                        commentsHTML += `
                            <div class="comment">
                                <strong>${comment.author}</strong>
                                <p>${comment.content}</p>
                            </div>
                        `;
                    });
                }

                article.innerHTML = `
                    <div class="post-header">
                        <h3>${post.title}</h3>
                        <span class="category-badge">${post.category}</span>
                    </div>

                    <p>${post.content}</p>

                    ${imageHTML}

                    <small>Posté par ${post.author}</small>

                    <div class="reaction-bar">
                        <button class="reaction-btn" onclick="reactToPost(${post.id}, 'like')">
                            👍 ${post.likes}
                        </button>

                        <button class="reaction-btn" onclick="reactToPost(${post.id}, 'dislike')">
                            👎 ${post.dislikes}
                        </button>
                    </div>

                    <div class="comments">
                        ${commentsHTML}
                    </div>

                    <form class="comment-form" data-post-id="${post.id}">
                        <input type="text" name="content" placeholder="Écrire un commentaire..." required>
                        <button type="submit">Commenter</button>
                    </form>
                `;

                postsList.appendChild(article);
            });

            document.querySelectorAll(".comment-form").forEach(form => {
                form.addEventListener("submit", async event => {
                    event.preventDefault();

                    const contentInput = form.querySelector("input");

                    const formData = new FormData();
                    formData.append("post_id", form.dataset.postId);
                    formData.append("content", contentInput.value);

                    const response = await fetch("/api/comments/add", {
                        method: "POST",
                        body: formData
                    });

                    if (response.ok) {
                        contentInput.value = "";
                        loadPosts();
                    } else {
                        alert("Vous devez être connecté pour commenter.");
                    }
                });
            });
        }

        function applyFilters() {
            const searchValue = searchInput ? searchInput.value.toLowerCase() : "";
            const selectedCategory = categoryFilter ? categoryFilter.value : "";

            const filteredPosts = posts.filter(post => {
                const matchesSearch =
                    post.title.toLowerCase().includes(searchValue) ||
                    post.content.toLowerCase().includes(searchValue) ||
                    post.author.toLowerCase().includes(searchValue);

                const matchesCategory =
                    selectedCategory === "" || post.category === selectedCategory;

                return matchesSearch && matchesCategory;
            });

            displayPosts(filteredPosts);
        }

        if (searchInput) {
            searchInput.addEventListener("input", applyFilters);
        }

        if (categoryFilter) {
            categoryFilter.addEventListener("change", applyFilters);
        }

        displayPosts(posts);

    } catch (error) {
        console.error(error);

        postsList.innerHTML = `
            <article class="post-card">
                <h3>Erreur</h3>
                <p>Impossible de charger les posts.</p>
            </article>
        `;
    }
}

async function reactToPost(postId, type) {
    const formData = new FormData();
    formData.append("post_id", postId);
    formData.append("type", type);

    const response = await fetch("/api/react", {
        method: "POST",
        body: formData
    });

    if (response.ok) {
        loadPosts();
    } else {
        alert("Vous devez être connecté pour liker ou disliker.");
    }
}