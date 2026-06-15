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
        const postsResponse = await fetch("/api/posts");
        const posts = await postsResponse.json();

        const statsResponse = await fetch("/api/stats");
        const stats = await statsResponse.json();

        let totalComments = 0;

        posts.forEach(post => {
            totalComments += post.comments ? post.comments.length : 0;
        });

        if (postsCount) {
            postsCount.textContent = posts.length;
        }

        if (usersCount) {
            usersCount.textContent = stats.users;
        }

        if (commentsCount) {
            commentsCount.textContent = totalComments;
        }

        function displayPosts(filteredPosts) {
            postsList.innerHTML = "";

            if (filteredPosts.length === 0) {
                postsList.innerHTML = `
                    <article class="post-card">
                        <h3>Aucun post trouvé</h3>
                        <p>Aucune discussion ne correspond à ta recherche.</p>
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
                        <img
                            class="post-image"
                            src="${post.image}"
                            alt="Image du post"
                        >
                    `;
                }

                let commentsHTML = "";

                if (post.comments && post.comments.length > 0) {
                    commentsHTML = `
                        <div class="comments-title">
                            Commentaires
                        </div>
                    `;

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
                        <div>
                            <h3>${post.title}</h3>
                            <small class="post-author">
                                Posté par ${post.author}
                            </small>
                        </div>

                        <span class="category-badge">
                            ${post.category}
                        </span>
                    </div>

                    <p class="post-content">
                        ${post.content}
                    </p>

                    ${imageHTML}

                    <div class="reaction-bar">
                        <button
                            class="reaction-btn"
                            onclick="reactToPost(${post.id}, 'like')"
                        >
                            👍 ${post.likes}
                        </button>

                        <button
                            class="reaction-btn"
                            onclick="reactToPost(${post.id}, 'dislike')"
                        >
                            👎 ${post.dislikes}
                        </button>
                    </div>

                    <div class="comments">
                        ${commentsHTML}
                    </div>

                    <form class="comment-form" data-post-id="${post.id}">
                        <input
                            type="text"
                            name="content"
                            placeholder="Écrire un commentaire..."
                            required
                        >

                        <button type="submit">
                            Commenter
                        </button>
                    </form>
                `;

                postsList.appendChild(article);
            });

            document.querySelectorAll(".comment-form").forEach(form => {
                form.addEventListener("submit", async event => {
                    event.preventDefault();

                    const input = form.querySelector("input");
                    const content = input.value.trim();

                    if (content === "") {
                        return;
                    }

                    const formData = new FormData();

                    formData.append("post_id", form.dataset.postId);
                    formData.append("content", content);

                    const response = await fetch("/api/comments/add", {
                        method: "POST",
                        body: formData
                    });

                    if (response.ok) {
                        input.value = "";
                        loadPosts();
                    } else {
                        alert("Vous devez être connecté pour commenter.");
                    }
                });
            });
        }

        function applyFilters() {
            const searchValue = searchInput
                ? searchInput.value.toLowerCase()
                : "";

            const selectedCategory = categoryFilter
                ? categoryFilter.value
                : "";

            const filteredPosts = posts.filter(post => {
                const title = post.title.toLowerCase();
                const content = post.content.toLowerCase();
                const author = post.author.toLowerCase();
                const category = post.category;

                const matchesSearch =
                    title.includes(searchValue) ||
                    content.includes(searchValue) ||
                    author.includes(searchValue);

                const matchesCategory =
                    selectedCategory === "" ||
                    category === selectedCategory;

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