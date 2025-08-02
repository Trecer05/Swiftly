import '../../domain/user/models/user.dart';

class ApiUser {
  final String id;
  final String? name;
  final String? lastName;
  final String? image;

  ApiUser({required this.id, this.name, this.lastName, this.image});

  User toDomain() {
    return User.create(
      id: id,
      name: name,
      lastName: lastName,
      image: image,
    );
  }
}
