import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent, act } from "@testing-library/react";
import { AuthProvider, useAuth } from "./AuthContext";

function Consumer() {
  const { user, isLoggedIn, loading, login, logout } = useAuth();
  return (
    <div>
      <span data-testid="status">
        {loading ? "loading" : isLoggedIn ? "in" : "out"}
      </span>
      <span data-testid="email">{user?.email ?? "none"}</span>
      <button onClick={() => login("idtok").catch(() => {})}>login</button>
      <button onClick={() => logout()}>logout</button>
    </div>
  );
}

function renderAuth() {
  return render(
    <AuthProvider>
      <Consumer />
    </AuthProvider>
  );
}

beforeEach(() => {
  localStorage.clear();
  vi.restoreAllMocks();
});

describe("AuthContext", () => {
  it("starts logged out when there is no stored jwt", async () => {
    renderAuth();
    await waitFor(() =>
      expect(screen.getByTestId("status").textContent).toBe("out")
    );
  });

  it("login stores the jwt and sets the user", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          token: "jwt123",
          user: { id: "u1", email: "a@b.com", credit: 10 },
        }),
      })
      .mockResolvedValue({
        ok: true,
        json: async () => ({ id: "u1", email: "a@b.com", credit: 10 }),
      });
    vi.stubGlobal("fetch", fetchMock);

    renderAuth();
    await act(async () => {
      fireEvent.click(screen.getByText("login"));
    });

    await waitFor(() =>
      expect(screen.getByTestId("email").textContent).toBe("a@b.com")
    );
    expect(localStorage.getItem("jwt")).toBe("jwt123");
  });

  it("login throws and stays logged out on a failed response", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: false }));
    renderAuth();
    await act(async () => {
      fireEvent.click(screen.getByText("login"));
    });
    await waitFor(() =>
      expect(screen.getByTestId("status").textContent).toBe("out")
    );
    expect(localStorage.getItem("jwt")).toBeNull();
  });

  it("restores the session from a stored jwt on mount", async () => {
    localStorage.setItem("jwt", "stored-jwt");
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ id: "u9", email: "stored@x.com", credit: 5 }),
      })
    );

    renderAuth();
    await waitFor(() =>
      expect(screen.getByTestId("email").textContent).toBe("stored@x.com")
    );
  });

  it("clears an expired/invalid stored jwt (401)", async () => {
    localStorage.setItem("jwt", "bad-jwt");
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: false }));

    renderAuth();
    await waitFor(() =>
      expect(screen.getByTestId("status").textContent).toBe("out")
    );
    expect(localStorage.getItem("jwt")).toBeNull();
  });

  it("logout clears the user and storage", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          token: "jwt123",
          user: { id: "u1", email: "a@b.com", credit: 10 },
        }),
      })
      .mockResolvedValue({
        ok: true,
        json: async () => ({ id: "u1", email: "a@b.com", credit: 10 }),
      });
    vi.stubGlobal("fetch", fetchMock);

    renderAuth();
    await act(async () => {
      fireEvent.click(screen.getByText("login"));
    });
    await waitFor(() =>
      expect(screen.getByTestId("email").textContent).toBe("a@b.com")
    );

    await act(async () => {
      fireEvent.click(screen.getByText("logout"));
    });
    await waitFor(() =>
      expect(screen.getByTestId("status").textContent).toBe("out")
    );
    expect(localStorage.getItem("jwt")).toBeNull();
  });
});
