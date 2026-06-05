import localFont from "next/font/local";

export const iannnnnDog = localFont({
  src: [
    {
      path: "./public/assets/fonts/iannnnn-DOG-Regular.ttf",
      weight: "400",
      style: "normal",
    },
    {
      path: "./public/assets/fonts/iannnnn-DOG-Bold.ttf",
      weight: "700",
      style: "normal",
    },
    {
      path: "./public/assets/fonts/iannnnn-DOG-Light.ttf",
      weight: "300",
      style: "normal",
    },
  ],
  variable: "--font-iannnnn-dog",
  display: "swap",
});