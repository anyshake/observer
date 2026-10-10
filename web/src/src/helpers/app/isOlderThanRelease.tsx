const parseVersion = (version: string) => {
    if (version === '<custom-version>') {
        return { numbers: [BigInt(0), BigInt(0), BigInt(0)], preRelease: '' };
    }

    const match = /^v(-?\d+)\.(-?\d+)\.(-?\d+)(?:-([\s\S]+))?$/.exec(version);
    if (!match) {
        return null;
    }

    return {
        numbers: [BigInt(match[1]), BigInt(match[2]), BigInt(match[3])],
        preRelease: match[4] ?? ''
    };
};

const isNumericIdentifier = (identifier: string) => /^[0-9]+$/.test(identifier);

const isLexicallyLess = (left: string, right: string) => {
    const encoder = new TextEncoder();
    const leftBytes = encoder.encode(left);
    const rightBytes = encoder.encode(right);
    for (let i = 0; i < Math.min(leftBytes.length, rightBytes.length); i++) {
        if (leftBytes[i] !== rightBytes[i]) {
            return leftBytes[i] < rightBytes[i];
        }
    }
    return leftBytes.length < rightBytes.length;
};

export const isOlderThanRelease = (current: string, release: string) => {
    const currentVersion = parseVersion(current);
    const releaseVersion = parseVersion(release);
    if (!currentVersion || !releaseVersion) {
        return false;
    }

    for (let i = 0; i < 3; i++) {
        if (currentVersion.numbers[i] !== releaseVersion.numbers[i]) {
            return currentVersion.numbers[i] < releaseVersion.numbers[i];
        }
    }

    if (currentVersion.preRelease && !releaseVersion.preRelease) {
        return true;
    }
    if (!currentVersion.preRelease && releaseVersion.preRelease) {
        return false;
    }

    const currentIdentifiers = currentVersion.preRelease.split('.');
    const releaseIdentifiers = releaseVersion.preRelease.split('.');
    for (let i = 0; i < Math.min(currentIdentifiers.length, releaseIdentifiers.length); i++) {
        const currentIdentifier = currentIdentifiers[i];
        const releaseIdentifier = releaseIdentifiers[i];
        if (currentIdentifier === releaseIdentifier) {
            continue;
        }

        const currentNumeric = isNumericIdentifier(currentIdentifier);
        const releaseNumeric = isNumericIdentifier(releaseIdentifier);
        if (currentNumeric !== releaseNumeric) {
            return currentNumeric;
        }
        if (currentNumeric && currentIdentifier.length !== releaseIdentifier.length) {
            return currentIdentifier.length < releaseIdentifier.length;
        }
        return isLexicallyLess(currentIdentifier, releaseIdentifier);
    }

    return currentIdentifiers.length < releaseIdentifiers.length;
};
