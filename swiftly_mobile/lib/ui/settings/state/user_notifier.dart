import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../domain/user/models/user.dart';
import 'user_state.dart';

class UserNotifier extends StateNotifier<UserState> {
  UserNotifier() : super(UserState.initial());

  void addUser(User user) {
    state = state.copyWith(users: [...state.users, user]);
  }

  void removeUser(String id) {
    state = state.copyWith(
      users: state.users.where((user) => user.id != id).toList(),
    );
  }

  void updateName(String userId, String newName) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newName) : user;
          }).toList(),
    );
  }

  void updateLastName(String userId, String newLastName) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newLastName) : user;
          }).toList(),
    );
  }

  void updateImage(String userId, String newImage) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newImage) : user;
          }).toList(),
    );
  }

  void updateRole(String userId, String newRole) {
    state = state.copyWith(
      users:
          state.users.map((user) {
            return user.id == userId ? user.copyWith(name: newRole) : user;
          }).toList(),
    );
  }

  void loadUsers() async {
    state = state.copyWith(isLoading: true);
    await Future.delayed(const Duration(seconds: 1));
    state = state.copyWith(
      users: [...state.users, ...mockUsers],
      isLoading: false,
    );
  }
}

final mockUsers = [
  User.create(
    id: '1',
    name: 'Павел',
    image:
        'https://65.mchs.gov.ru/uploads/resize_cache/news/2021-08-25/pravila-povedeniya-pri-vstreche-s-medvedem_1629847892112633638__800x800.jpg',
  ),
  User.create(
    id: '2',
    name: 'Добрыня',
    lastName: 'Никитич',
  ),
  User.create(
    id: '3',
    name: 'Алеша',
    lastName: 'Попович',
  ),
  User.create(
    id: '4',
    name: 'Анастасия',
    lastName: 'Петровна',
    image:
        'https://cs13.pikabu.ru/post_img/big/2024/03/06/5/1709705621175092550.png',
  ),
  User.create(
    id: '5',
    name: 'Алексей',
    lastName: 'Семеныч',
  ),
];
