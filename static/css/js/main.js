document.addEventListener("DOMContentLoaded", () => {
    console.log("Nexora JS chargé");

    loadPosts();
});

async function loadPosts() {
    const postsList = document.getElementById("postsList");
    const searchInput = document.getElementById("searchInput");

    const postsCount = document.getElementById("postsCount");
    const usersCount = document.getElementById("usersCount");
    const commentsCount = document.getElementById("commentsCount");

    try {
        const response = await fetch("/api/posts");
        const posts = await response.json();

        postsList.innerHTML = "";

        animateCounter(postsCount, posts.length);
        animateCounter(usersCount, 1);
        animateCounter(commentsCount, 0);

        if (posts.length === 0) {
            postsList.innerHTML = `
                <article class="post-card">
                    <h3>Aucun post</h3>
                    <p>Soyez le premier à publier sur Nexora.</p>
                </article>
            `;
            return;
        }

        posts.forEach(post => {
            const article = document.createElement("article");

            article.className = "post-card";

            article.dataset.title = post.title;

            article.innerHTML = `
                <h3>${post.title}</h3>
                <p>${post.content}</p>

                <small>
                    Posté par ${post.author}
                    • ${post.category}
                </small>

                <form class="comment-form" data-post-id="${post.id}">
                    <input
                        type="text"
                        name="content"
                        placeholder="Écrire un commentaire..."
                    >

                    <button type="submit">
                        Commenter
                    </button>
                </form>
            `;

            postsList.appendChild(article);
        });

        document.querySelectorAll(".comment-form").forEach(form => {
            form.addEventListener("submit", async (event) => {
                event.preventDefault();

                const postId = form.dataset.postId;

                const content =
                    form.querySelector("input").value;

                const formData = new FormData();

                formData.append("post_id", postId);
                formData.append("content", content);

                const response = await fetch(
                    "/api/comments/add",
                    {
                        method: "POST",
                        body: formData,
                    }
                );

                if (response.ok) {
                    form.querySelector("input").value = "";

                    alert("Commentaire ajouté !");
                } else {
                    alert(
                        "Vous devez être connecté pour commenter."
                    );
                }
            });
        });

        if (searchInput) {
            searchInput.addEventListener("input", () => {
                const searchValue =
                    searchInput.value.toLowerCase();

                document
                    .querySelectorAll(".post-card")
                    .forEach(post => {
                        const title =
                            post.dataset.title.toLowerCase();

                        const text =
                            post.textContent.toLowerCase();

                        if (
                            title.includes(searchValue) ||
                            text.includes(searchValue)
                        ) {
                            post.style.display = "block";
                        } else {
                            post.style.display = "none";
                        }
                    });
            });
        }
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

function animateCounter(element, target) {
    if (!element) return;

    let count = 0;

    const interval = setInterval(() => {
        element.textContent = count;

        if (count >= target) {
            clearInterval(interval);
        }

        count++;
    }, 120);
}