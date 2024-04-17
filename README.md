# Networking

The goal of this project to intimately learn each of 7 layers of the [OSI model](https://en.wikipedia.org/wiki/OSI_model) by simulating and implementing a simplified version of the logic of each layer.

### Physical layer

- I didn't want to take anything for granted in this layer
- The biggest limitation I imposed on myself was that the only way two _devices_ could communicate was via a shared `Cable` object that was either `Postive` or `Negative`
  - The cable being on or off brings along with it no timing information
  - To solve for this problem, I used what is known as manchester encoding for which transitions between states carries information
  - I then implmeneted an algorithm for guessing the clock rate of the incoming the bit stream that was able to recover from collisions and timing issues
    - Credit to [victornpb/manch_decode](https://github.com/victornpb/manch_decode) for some inspiration of how to implement
