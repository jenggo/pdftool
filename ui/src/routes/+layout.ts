import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import type { LayoutLoad } from './$types';

export const prerender = true;
export const load: LayoutLoad = async ({ url }) => {
  // Skip auth check for login page
  if (url.pathname === '/login') {
    return {};
  }

  // Check if user is authenticated
  if (browser) {
    try {
      const response = await fetch('/check-auth');
      if (!response.ok) {
        goto('/login');
      }
    } catch (error) {
      goto('/login');
    }
  }

  return {};
};
