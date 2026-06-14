
import type { Component } from "svelte";

export type Page = {
  id: string;
  title: string;
  icon: Component;
  section: "header" | "content" | "footer";
  component: Component | null;
  onInit: () => void | null;
};

export type ErrorPromise = Promise<string | null>
export type ErrorPromiseArr = Promise<string[] | null>
