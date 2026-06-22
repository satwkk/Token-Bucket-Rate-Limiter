import requests
import threading
import time

URL = "http://localhost:8080/"
HEADERS = {"X-Rate-Limit-ID": "127.0.0.1"}
WORKERS = 25

lock = threading.Lock()
blocked_at = None
was_blocked = False
stats = {"total": 0, "blocked": 0, "allowed": 0}


def handle_response(status_code: int) -> None:
    global blocked_at, was_blocked

    with lock:
        stats["total"] += 1
        now = time.time()

        if status_code == 429:
            stats["blocked"] += 1
            if not was_blocked:
                was_blocked = True
                blocked_at = now
                print(f"[#{stats['total']}] BLOCKED at {time.strftime('%H:%M:%S')}")
            return

        stats["allowed"] += 1
        if was_blocked and blocked_at is not None:
            elapsed = now - blocked_at
            print(
                f"[#{stats['total']}] ALLOWED after {elapsed:.2f}s "
                f"(status {status_code})"
            )
            was_blocked = False
            blocked_at = None
        elif stats["allowed"] <= 5 or stats["allowed"] % 50 == 0:
            print(f"[#{stats['total']}] OK (status {status_code})")


def worker(session: requests.Session) -> None:
    while True:
        try:
            response = session.get(URL, headers=HEADERS, timeout=5)
            handle_response(response.status_code)
        except requests.RequestException as exc:
            with lock:
                stats["total"] += 1
            print(f"[#{stats['total']}] Request error: {exc}")


def main() -> None:
    print(f"Hammering {URL} with {WORKERS} workers (Ctrl+C to stop)\n")

    threads = []
    for _ in range(WORKERS):
        session = requests.Session()
        thread = threading.Thread(target=worker, args=(session,), daemon=True)
        thread.start()
        threads.append(thread)

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        with lock:
            print(
                f"\nStopped. total={stats['total']} "
                f"allowed={stats['allowed']} blocked={stats['blocked']}"
            )


if __name__ == "__main__":
    main()
