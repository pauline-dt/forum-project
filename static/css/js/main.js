document.addEventListener("DOMContentLoaded", () => {
    console.log("Nexora JS chargé");

    const posts = document.querySelectorAll(".post-card");
    const searchInput = document.getElementById("searchInput");

    const postsCount = document.getElementById("postsCount");
    const usersCount = document.getElementById("usersCount");
    const commentsCount = document.getElementById("commentsCount");

    animateCounter(postsCount, 3);
    animateCounter(usersCount, 1);
    animateCounter(commentsCount, 0);

    if (searchInput) {
        searchInput.addEventListener("input", () => {
            const searchValue = searchInput.value.toLowerCase();

            posts.forEach(post => {
                const title = post.dataset.title.toLowerCase();
                const text = post.textContent.toLowerCase();

                if (title.includes(searchValue) || text.includes(searchValue)) {
                    post.style.display = "block";
                } else {
                    post.style.display = "none";
                }
            });
        });
    }
});

function animateCounter(element, target) {
    let count = 0;

    const interval = setInterval(() => {
        element.textContent = count;

        if (count >= target) {
            clearInterval(interval);
        }

        count++;
    }, 120);
}