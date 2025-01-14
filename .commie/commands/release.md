# Create a good release tag.

Step 1. Examine all the git tags.

Step 2. Ask a user for the release type (showing the new version number with each bump type):
 1. patch - VERSION_FOR_PATCH
 2. minor - VERSION_FOR_MINOR
 3. major - VERSION_FOR_MAJOR
 4. pre-release - VERSION_FOR_PRERELEASE
 5. other (please specify in the comment)

Create a new **unsigned** tag according to the desired bump type.

Create the tag. Report about it.

Do `git push && git push --tags`