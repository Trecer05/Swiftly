import 'package:freezed_annotation/freezed_annotation.dart';

part 'user.freezed.dart';

@freezed
class User with _$User {
  const factory User({
    required String id,
    required String name,
    required String? lastName,
    required String? image,
  }) = _User;

  factory User.create({
    required String id,
    String? name,
    String? lastName,
    String? image,
  }) => User(
    id: id,
    name: name?.trim().isEmpty ?? true ? 'Новый пользователь' : name!,
    lastName: lastName,
    image: image,
  );
}
