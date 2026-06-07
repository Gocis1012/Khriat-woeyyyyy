import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { GuestProvider, useGuest } from "./GuestContext";

// Control the auth state the GuestProvider sees.
const { authState } = vi.hoisted(() => ({
  authState: {
    value: {
      user: null as null | { credit: number },
      token: null as string | null,
      isLoggedIn: false,
      loading: false,
    },
  },
}));

vi.mock("./AuthContext", () => ({
  useAuth: () => authState.value,
}));

function Consumer() {
  const { credit } = useGuest();
  return <span data-testid="credit">{credit === null ? "null" : credit}</span>;
}

function renderGuest() {
  return render(
    <GuestProvider>
      <Consumer />
    </GuestProvider>
  );
}

beforeEach(() => {
  vi.restoreAllMocks();
  authState.value = { user: null, token: null, isLoggedIn: false, loading: false };
});

describe("GuestContext", () => {
  it("fetches guest credit from /guest/status when not logged in", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ credit: 6 }),
    });
    vi.stubGlobal("fetch", fetchMock);

    renderGuest();
    await waitFor(() =>
      expect(screen.getByTestId("credit").textContent).toBe("6")
    );
    expect(fetchMock.mock.calls[0][0]).toContain("/guest/status");
  });

  it("uses the logged-in user's credit and the /user/me endpoint", async () => {
    authState.value = {
      user: { credit: 10 },
      token: "jwt",
      isLoggedIn: true,
      loading: false,
    };
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ credit: 10 }),
    });
    vi.stubGlobal("fetch", fetchMock);

    renderGuest();
    await waitFor(() =>
      expect(screen.getByTestId("credit").textContent).toBe("10")
    );
    expect(fetchMock.mock.calls[0][0]).toContain("/api/v1/user/me");
  });

  it("silently ignores fetch errors (credit stays null)", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("network")));
    renderGuest();
    // Give the effect a tick; credit should remain null.
    await new Promise((r) => setTimeout(r, 20));
    expect(screen.getByTestId("credit").textContent).toBe("null");
  });
});
